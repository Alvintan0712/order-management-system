package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"example.com/oms/common"
	"example.com/oms/common/broker/consumer"
	"example.com/oms/common/broker/topic"
	"example.com/oms/common/discovery"
	"example.com/oms/common/discovery/consul"
	"example.com/oms/common/schemaregistry/deserializer"
	"example.com/oms/common/schemaregistry/srclient"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
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

	schemaRegistryURL = common.EnvString("SCHEMA_REGISTRY_URL", "http://localhost:8081")

	debug = common.EnvString("DEBUG", "false") == "true"
)

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

	confluentSRClient, err := srclient.NewConfluentSRClient(schemaRegistryURL)
	if err != nil {
		log.Fatal(err)
	}

	confluentDeserializer, err := deserializer.NewConfluentDeserializer(confluentSRClient, serde.ValueSerde)
	if err != nil {
		log.Fatal(err)
	}

	topics := []string{
		topic.MenuCreated,
	}
	kafkaConsumer, err := consumer.NewKafkaConsumer(kafkaBrokers, kafkaConsumerGroupId, topics)
	if err != nil {
		log.Fatal(err)
	}

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
	consumerHandler := NewConsumerHandler(confluentDeserializer, service)
	if err := kafkaConsumer.Start(consumerHandler.Handlers); err != nil {
		log.Fatal(err)
	}

	log.Printf("Stock service %s started at %s:%s\n", serviceId, serviceHost, servicePort)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err.Error())
	}
}
