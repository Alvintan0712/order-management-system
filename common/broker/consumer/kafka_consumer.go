package consumer

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	fmt.Println("Kafka consumer subscribing topics:", topics)

	kafkaConsumer := &KafkaConsumer{
		client: consumer,
	}

	return kafkaConsumer, nil
}

func (c *KafkaConsumer) Start(handlers map[string]EventHandler) error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	run := true
	for run {
		select {
		case <-signalChan:
			fmt.Println("Consumer shutting down...")
			run = false
		default:
			msg, err := c.client.ReadMessage(-1)
			if err != nil {
				log.Printf("Kafka consume error: %v", err)
				break
			}

			EventMessage := EventMessage{
				Topic:     *msg.TopicPartition.Topic,
				Partition: msg.TopicPartition.Partition,
				Offset:    int64(msg.TopicPartition.Offset),
				Key:       msg.Key,
				Value:     msg.Value,
				Timestamp: msg.Timestamp,
			}

			err = handlers[*msg.TopicPartition.Topic](EventMessage)
			if err != nil {
				log.Println(err)
				break
			}

			_, err = c.client.CommitMessage(msg)
			if err != nil {
				log.Printf("message commit error: %v", err)
			}
		}
	}

	return nil
}

func (c *KafkaConsumer) Close() error {
	return c.client.Close()
}
