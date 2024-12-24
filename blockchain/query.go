package blockchain

import (
	"encoding/json"
	"fmt"
	"time"
)

type Metadata struct {
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
}

func (c *Client) QueryMetadata(fileName string) (*Metadata, error) {
	start := time.Now()

	response, err := c.ChannelClient.EvaluateTransaction("QueryMetadata", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to query chaincode: %v", err)
	}

	latency := time.Since(start)
	fmt.Printf("Chaincode query latency: %s\n", latency)

	var metadata Metadata
	err = json.Unmarshal(response, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %v", err)
	}

	return &metadata, nil
}
