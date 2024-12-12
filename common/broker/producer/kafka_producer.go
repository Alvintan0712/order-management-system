package producer

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	client *kafka.Producer
}

func NewKafkaProducer(brokers, id string) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"client.id":         id,
		"acks":              "all",
		"retries":           5,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	kafkaProducer := &KafkaProducer{
		client: producer,
	}
	fmt.Println("Kafka producer started...")

	return kafkaProducer, nil
}

func (p *KafkaProducer) Produce(topic string, key, value []byte) error {
	deliveryChan := make(chan kafka.Event)

	err := p.client.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   key,
		Value: value,
	}, deliveryChan)
	if err != nil {
		return err
	}

	event := <-deliveryChan
	msg := event.(*kafka.Message)

	if msg.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %v", msg.TopicPartition.Error)
	}

	close(deliveryChan)

	return nil
}

func (p *KafkaProducer) Close() error {
	p.client.Close()
	return nil
}
