package main

import (
	"log"
	"proteomic/chaincode"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func main() {
	metadataChainCode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := metadataChainCode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}
