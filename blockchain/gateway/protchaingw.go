package gateway

import (
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"protchain/internal/config"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type FabricClient struct {
	Gateway      *client.Gateway
	Contract     *client.Contract
	Chaincode    string
	Channel      string
	PeerEndpoint string
	MspID        string
	CryptoPath   string
	TlsCertPath  string
	CertPath     string
	KeyPath      string
}

func NewFabricClient(cfg *config.Config) *FabricClient {
	// retrieve key path, as this is generated on network startup
	kp := path.Join(cfg.CryptoPath, "users/User1@org1.example.com/msp/keystore/")
	kp, err := findPrivateKey(kp)
	fmt.Println("keystore is", kp)
	fmt.Println("end key store")
	if err != nil {
		log.Fatalf("failed to discover secret key -> ", err.Error())
	}

	fabricClient := &FabricClient{
		Chaincode:    cfg.ChainCode,
		Channel:      cfg.Channel,
		PeerEndpoint: cfg.PeerEndpoint,
		MspID:        cfg.MSPID,
		CryptoPath:   cfg.CryptoPath,
		TlsCertPath:  path.Join(cfg.CryptoPath, "peers/peer0.org1.example.com/tls/ca.crt"),
		CertPath:     path.Join(cfg.CryptoPath, "users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"),
		KeyPath:      kp,
	}

	grpcConn, err := fabricClient.newGrpcConnection()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	id := fabricClient.newIdentity()
	sign := fabricClient.newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(grpcConn),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("failed to create gateway connection: %w", err)
		return nil
	}

	network := gw.GetNetwork(cfg.Channel)
	contract := network.GetContract(cfg.ChainCode)

	fabricClient.Gateway = gw
	fabricClient.Contract = contract

	return fabricClient
}

func findPrivateKey(path string) (string, error) {
	files, err := filepath.Glob(filepath.Join(path, "*_sk"))
	if err != nil || len(files) == 0 {
		return "", fmt.Errorf("private key not found in %s", path)
	}
	return files[0], nil
}

func (fc *FabricClient) newGrpcConnection() (*grpc.ClientConn, error) {
	certificatePEM, err := os.ReadFile(fc.TlsCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read TLS certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(certificatePEM)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, "")

	connection, err := grpc.Dial(fc.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("failed to establish gRPC connection: %w", err)
	}

	return connection, nil
}

func (fc *FabricClient) newIdentity() *identity.X509Identity {
	certificatePEM, err := os.ReadFile(fc.CertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(fmt.Errorf("failed to parse certificate: %w", err))
	}

	id, err := identity.NewX509Identity(fc.MspID, certificate)
	if err != nil {
		panic(fmt.Errorf("failed to create identity: %w", err))
	}

	return id
}

func (fc *FabricClient) newSign() identity.Sign {
	privateKeyPEM, err := os.ReadFile(fc.KeyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(fmt.Errorf("failed to parse private key: %w", err))
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(fmt.Errorf("failed to create signer: %w", err))
	}

	return sign
}

func (fc *FabricClient) SubmitTransaction(function string, args ...string) (string, error) {
	fmt.Printf("--> Submit Transaction: %s\n", function)

	result, err := fc.Contract.SubmitTransaction(function, args...)
	if err != nil {
		return "", fmt.Errorf("failed to submit transaction: %w", err)
	}

	fmt.Printf("*** Transaction committed successfully\n")
	return string(result), nil
}

func (fc *FabricClient) EvaluateTransaction(function string, args ...string) (string, error) {
	fmt.Printf("--> Evaluate Transaction: %s\n", function)

	result, err := fc.Contract.EvaluateTransaction(function, args...)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return string(result), nil
}
