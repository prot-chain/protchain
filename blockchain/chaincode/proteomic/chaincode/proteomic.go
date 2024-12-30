package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// ProteinMetadata refers to the metadata
// persisted to the blockchain for a protein
type ProteinMetadata struct {
	// Hash refers to the computed hash of the
	// protein data
	Hash      string `json:"protein_hash"`
	ProteinID string `json:"protein_id"`
	FileUrl   string `json:"file_Url"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	metadatas := []ProteinMetadata{
		{Hash: "testHash", ProteinID: "TES1", FileUrl: "https://test.com"},
		{Hash: "testHashTwo", ProteinID: "TES2", FileUrl: "https://test2.com"},
		{Hash: "testHashThree", ProteinID: "TES3", FileUrl: "https://test3.com"},
	}

	for _, metadata := range metadatas {
		metaJSON, err := json.Marshal(metadata)
		if err != nil {
			return err
		}

		if err := ctx.GetStub().PutState(metadata.ProteinID, metaJSON); err != nil {
			return errors.Wrap(err, "failed to put metadata to world state")
		}
	}

	return nil
}

func (s *SmartContract) StoreMetadata(ctx contractapi.TransactionContextInterface, proteinHash, proteinId, fileUrl string) error {
	metadata := ProteinMetadata{
		Hash:      proteinHash,
		ProteinID: proteinId,
		FileUrl:   fileUrl,
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %v", err)
	}

	return ctx.GetStub().PutState(proteinId, metadataJSON)
}

func (s *SmartContract) QueryMetadata(ctx contractapi.TransactionContextInterface, proteinId string) (*ProteinMetadata, error) {
	metadataBytes, err := ctx.GetStub().GetState(proteinId)
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %v", err)
	}

	if metadataBytes == nil {
		return nil, fmt.Errorf("metadata not found for protein: %s", proteinId)
	}

	var metadata ProteinMetadata
	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	return &metadata, nil
}

func (s *SmartContract) MetadataExists(ctx contractapi.TransactionContextInterface, proteinId string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(proteinId)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func (s *SmartContract) UpdateMetadata(ctx contractapi.TransactionContextInterface, proteinHash, proteinId, fileUrl string) error {
	exists, err := s.MetadataExists(ctx, proteinId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("No metadata exists for protein " + proteinId)
	}
	metadata := ProteinMetadata{
		Hash:      proteinHash,
		ProteinID: proteinId,
		FileUrl:   fileUrl,
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %v", err)
	}

	return ctx.GetStub().PutState(proteinId, metadataJSON)
}
