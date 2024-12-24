package proteomic_data

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Metadata struct {
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
}

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) StoreMetadata(ctx contractapi.TransactionContextInterface, fileName, fileURL string) error {
	metadata := Metadata{
		FileName: fileName,
		FileURL:  fileURL,
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %v", err)
	}

	return ctx.GetStub().PutState(fileName, metadataJSON)
}

func (s *SmartContract) QueryMetadata(ctx contractapi.TransactionContextInterface, fileName string) (*Metadata, error) {
	metadataBytes, err := ctx.GetStub().GetState(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %v", err)
	}

	if metadataBytes == nil {
		return nil, fmt.Errorf("metadata not found for file: %s", fileName)
	}

	var metadata Metadata
	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	return &metadata, nil
}
