package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"example.com/oms/common"
	"example.com/oms/common/discovery"
	"example.com/oms/common/discovery/consul"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	consulAddr  = common.EnvString("CONSUL_ADDR", "127.0.0.1:8500")
	serviceHost = common.EnvString("SERVICE_HOST", "127.0.0.1")
	servicePort = common.EnvString("SERVICE_PORT", "8083")
	serviceName = common.EnvString("SERVICE_NAME", "stock-service")

	kafkaBrokers         = common.EnvString("KAFKA_BROKERS", "localhost:9092")
	kafkaConsumerGroupId = common.EnvString("KAFKA_CONSUMER_GROUP_ID", "stock-consumer-group")
	kafkaTopics          = common.EnvString("KAFKA_TOPICS", "")

	debug = common.EnvString("DEBUG", "false") == "true"
)

func initKafkaAdminClient() (*kafka.AdminClient, error) {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBrokers,
	})
	if err != nil {
		log.Printf("failed to create admin client: %s\n", err)
		return nil, err
	}

	topics := strings.Split(kafkaTopics, ",")
	if err := common.SetupKafkaTopics(adminClient, topics); err != nil {
		return nil, err
	}

	return adminClient, nil
}

func startKafkaConsumer() (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaBrokers,
		"group.id":           kafkaConsumerGroupId,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		log.Printf("failed to create Kafka consumer: %v\n", err)
		return nil, err
	}

	topics := strings.Split(kafkaTopics, ",")
	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Printf("failed to subscribe to topics: %v\n", err)
		return nil, err
	}

	fmt.Println("Kafka consumer started...")
	go func() {
		for {
			msg, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				log.Printf("failed to read message: %v\n", err)
				continue
			}

			// consume event (replace it to other function)
			fmt.Printf("Consumed event from topic: %s: key = %-10s value = %s\n", *msg.TopicPartition.Topic, string(msg.Key), string(msg.Value))
			_, err = consumer.CommitMessage(msg)
			if err != nil {
				log.Printf("failed to commit message: %v\n", err)
			}
		}
	}()

	return consumer, nil
}

func main() {
	ctx := context.Background()

	if debug {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stderr))
	}

	registry, err := consul.NewRegistry(consulAddr, serviceHost, servicePort, serviceName)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v\n", err)
	}

	serviceId := discovery.GenerateInstanceId(serviceName)
	for {
		err = registry.Register(ctx, serviceId, serviceName, serviceHost, servicePort)
		if err != nil {
			log.Printf("Error registering service with Consul: %v", err)
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	defer registry.Deregister(ctx, serviceId, serviceName)

	go func() {
		for {
			if err := registry.HealthCheck(serviceId, serviceName); err != nil {
				log.Fatalf("Failed to health check: %v", err)
			}
			time.Sleep(time.Second)
		}
	}()

	adminClient, err := initKafkaAdminClient()
	if err != nil {
		log.Fatalf("Kafka initial error: %v\n", err)
	}
	defer adminClient.Close()

	// start kafka consumer
	consumer, err := startKafkaConsumer()
	if err != nil {
		log.Fatalf("Kafka consumer error: %v", err)
	}
	defer consumer.Close()

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", serviceHost, servicePort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listen.Close()

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	repository, err := NewRepository()
	if err != nil {
		log.Fatalf("repository create failed: %v\n", err)
	}

	service, err := NewService(ctx, repository, registry)
	if err != nil {
		log.Fatalf("service create failed: %v\n", err)
	}

	NewGRPCHandler(grpcServer, service)

	log.Printf("Stock service %s started at %s:%s\n", serviceId, serviceHost, servicePort)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err.Error())
	}
}
