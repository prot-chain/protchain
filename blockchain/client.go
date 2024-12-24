package blockchain

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type Client struct {
	ChannelClient *channel.Client
	Gateway       *gateway.Gateway
}

func NewClient(configPath string, walletPath string, user string, channelName string) (*Client, error) {
	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %v", err)
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(configPath),
		gateway.WithIdentity(wallet, user),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gateway: %v", err)
	}

	network := gw.GetNetwork(channelName)
	if network == nil {
		return nil, fmt.Errorf("failed to connect to network: %s", channelName)
	}

	channelClient, err := network.GetContract("mychaincode")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to contract: %v", err)
	}

	return &Client{
		ChannelClient: channelClient,
		Gateway:       gw,
	}, nil
}
