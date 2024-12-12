package srclient

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
)

type ConfluentSRClient struct {
	Client schemaregistry.Client
}

func NewConfluentSRClient(url string) (*ConfluentSRClient, error) {
	config := schemaregistry.NewConfig(url)
	client, err := schemaregistry.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create schema registry: %v", err)
	}

	confluentSRClient := &ConfluentSRClient{
		Client: client,
	}

	return confluentSRClient, nil
}

func (c *ConfluentSRClient) Close() error {
	return c.Client.Close()
}
