package common

import (
	"context"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func CheckTopicExists(adminClient *kafka.AdminClient, topic string) bool {
	for {
		metadata, err := adminClient.GetMetadata(&topic, false, 5000)
		if err != nil {
			log.Printf("Error fetching metadata: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		for t := range metadata.Topics {
			if t == topic {
				return true
			}
		}

		break
	}

	return false
}

func SetupKafkaTopics(adminClient *kafka.AdminClient, topics []string) error {
	topicSpecifications := []kafka.TopicSpecification{}
	for _, topic := range topics {
		if CheckTopicExists(adminClient, topic) {
			log.Println("topic", topic, "exists")
			continue
		}

		topicSpecifications = append(topicSpecifications, kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 3,
		})
	}

	log.Println("topics specification:", topicSpecifications)

	if len(topicSpecifications) > 0 {
		_, err := adminClient.CreateTopics(context.Background(), topicSpecifications)
		if err != nil {
			return err
		}
	}

	return nil
}
