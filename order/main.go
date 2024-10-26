package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"example.com/oms/common"
	"example.com/oms/common/discovery"
	"example.com/oms/common/discovery/consul"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

var (
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	serviceHost = common.EnvString("SERVICE_HOST", "localhost")
	servicePort = common.EnvString("SERVICE_PORT", "8081")
	serviceName = common.EnvString("SERVICE_NAME", "order-service")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceHost, servicePort, serviceName)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	serviceId := discovery.GenerateInstanceId(serviceName)
	err = registry.Register(context.Background(), serviceId, serviceName, serviceHost, servicePort)
	if err != nil {
		log.Fatalf("Error registering service with Consul: %v", err)
		panic(err)
	}
	defer registry.Deregister(context.Background(), serviceId, serviceName)

	go func() {
		for {
			if err := registry.HealthCheck(serviceId, serviceName); err != nil {
				log.Fatalf("Failed to health check: %v", err)
			}
			time.Sleep(time.Second)
		}
	}()

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	repository := NewRepository()
	service := NewService(repository)

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", serviceHost, servicePort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listen.Close()

	NewGRPCHandler(grpcServer, service)

	log.Printf("Order service %s started at %s:%s\n", serviceId, serviceHost, servicePort)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err.Error())
	}
}
