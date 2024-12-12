package consumer

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumer struct {
	client *kafka.Consumer
}

func NewKafkaConsumer(brokers, id string, topics []string) (*KafkaConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  brokers,
		"group.id":           id,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %v", err)
	}

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe Kafka topics: %v", err)
	}

	kafkaConsumer := &KafkaConsumer{
		client: consumer,
	}
	fmt.Println("Kafka consumer started...")

	return kafkaConsumer, nil
}

func (c *KafkaConsumer) Consume() (*ConsumedMessage, error) {
	msg, err := c.client.ReadMessage(-1)
	if err != nil {
		return nil, err
	}

	consumedMessage := &ConsumedMessage{
		Topic:     *msg.TopicPartition.Topic,
		Partition: msg.TopicPartition.Partition,
		Offset:    int64(msg.TopicPartition.Offset),
		Key:       msg.Key,
		Value:     msg.Value,
	}

	return consumedMessage, nil
}

func (c *KafkaConsumer) Close() error {
	return c.client.Close()
}
