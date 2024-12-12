package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"example.com/oms/common"
	"example.com/oms/common/broker/producer"
	"example.com/oms/common/discovery"
	"example.com/oms/common/discovery/consul"
	"example.com/oms/common/schemaregistry/serializer"
	"example.com/oms/common/schemaregistry/srclient"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	consulAddr  = common.EnvString("CONSUL_ADDR", "127.0.0.1:8500")
	serviceHost = common.EnvString("SERVICE_HOST", "127.0.0.1")
	servicePort = common.EnvString("SERVICE_PORT", "8082")
	serviceName = common.EnvString("SERVICE_NAME", "menu-service")

	kafkaBrokers = common.EnvString("KAFKA_BROKERS", "localhost:9092")

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
				log.Fatalf("Failed to health check: %v\n", err)
			}
			time.Sleep(time.Second)
		}
	}()

	kafkaProducer, err := producer.NewKafkaProducer(kafkaBrokers, "menu-service")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer kafkaProducer.Close()

	srclient, err := srclient.NewConfluentSRClient(schemaRegistryURL)
	if err != nil {
		log.Fatal(err)
	}

	confluentSerializer, err := serializer.NewConfluentSerializer(srclient, serde.ValueSerde)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Start listen %s:%s\n", serviceHost, servicePort)
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

	service := NewService(repository, kafkaProducer, confluentSerializer)
	NewGRPCHandler(grpcServer, service)

	log.Printf("Menu service %s started at %s:%s\n", serviceId, serviceHost, servicePort)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err.Error())
	}
}
