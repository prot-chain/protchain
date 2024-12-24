package blockchain

import (
	"fmt"
	"time"
)

func (c *Client) StoreMetadata(fileName, fileURL string) error {
	start := time.Now()

	_, err := c.ChannelClient.SubmitTransaction("StoreMetadata", fileName, fileURL)
	if err != nil {
		return fmt.Errorf("failed to invoke chaincode: %v", err)
	}

	latency := time.Since(start)
	fmt.Printf("Chaincode invocation latency: %s\n", latency)
	return nil
}
