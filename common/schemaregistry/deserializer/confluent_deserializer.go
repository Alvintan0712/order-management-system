package deserializer

import (
	"fmt"

	"example.com/oms/common/schemaregistry/srclient"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avrov2"
)

type ConfluentDeserializer struct {
	SRClient     *srclient.ConfluentSRClient
	Deserializer *avrov2.Deserializer
}

func NewConfluentDeserializer(srclient *srclient.ConfluentSRClient, serdeType int) (*ConfluentDeserializer, error) {
	config := avrov2.NewDeserializerConfig()
	deser, err := avrov2.NewDeserializer(srclient.Client, serdeType, config)
	if err != nil {
		return nil, fmt.Errorf("create confluent deserializer failed: %v", err)
	}

	confluentDeserializer := &ConfluentDeserializer{
		SRClient:     srclient,
		Deserializer: deser,
	}

	return confluentDeserializer, nil
}

func (deser *ConfluentDeserializer) Deserialize(topic string, payload []byte, value interface{}) error {
	err := deser.Deserializer.DeserializeInto(topic, payload, value)
	if err != nil {
		return fmt.Errorf("failed to deserialize payload: %v", err)
	}

	return nil
}
