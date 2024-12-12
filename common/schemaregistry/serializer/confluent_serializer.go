package serializer

import (
	"fmt"

	"example.com/oms/common/schemaregistry/srclient"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avrov2"
)

type ConfluentSerializer struct {
	SRClient   *srclient.ConfluentSRClient
	Serializer *avrov2.Serializer
}

func NewConfluentSerializer(srclient *srclient.ConfluentSRClient, serdeType int) (*ConfluentSerializer, error) {
	config := avrov2.NewSerializerConfig()
	ser, err := avrov2.NewSerializer(srclient.Client, serdeType, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create confluent serializer: %v", err)
	}

	confluentSerializer := &ConfluentSerializer{
		SRClient:   srclient,
		Serializer: ser,
	}

	return confluentSerializer, nil
}

func (ser *ConfluentSerializer) Serialize(topic string, value interface{}) ([]byte, error) {
	payload, err := ser.Serializer.Serialize(topic, value)
	if err != nil {
		return nil, fmt.Errorf("failed to serialized payload: %v", err)
	}

	return payload, nil
}
