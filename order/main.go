package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"example.com/oms/common"
	capi "github.com/hashicorp/consul/api"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	serviceHost = common.EnvString("SERVICE_HOST", "localhost")
	servicePort = common.EnvString("SERVICE_PORT", "8081")
	serviceId   = common.EnvString("SERVICE_ID", "order-service-id")
	serviceName = common.EnvString("SERVICE_NAME", "order-service")
)

func main() {
	config := capi.DefaultConfig()
	config.Address = consulAddr
	consulClient, err := capi.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", serviceHost, servicePort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listen.Close()

	portNumber, err := strconv.Atoi(servicePort)
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	var opts []grpc.ServerOption
	agentServiceRegistration := &capi.AgentServiceRegistration{
		ID:      serviceId,
		Name:    serviceName,
		Address: serviceHost,
		Port:    portNumber,
		Check: &capi.AgentServiceCheck{
			GRPC:       fmt.Sprintf("%s:%s", serviceHost, servicePort),
			Interval:   "10s",
			Timeout:    "30s",
			GRPCUseTLS: false,
		},
	}

	err = consulClient.Agent().ServiceRegister(agentServiceRegistration)
	if err != nil {
		log.Fatalf("Error registering service with Consul: %v", err)
	}

	grpcServer := grpc.NewServer(opts...)
	healthServer := health.NewServer()
	repository := NewRepository()
	service := NewService(repository)

	NewGRPCHandler(grpcServer, service)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	service.CreateOrder(context.Background())

	log.Println("Order server started at", fmt.Sprintf("%s:%s", serviceHost, servicePort))

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err.Error())
	}

	err = consulClient.Agent().ServiceDeregister(serviceId)
	if err != nil {
		log.Printf("Error deregistering service: %v", err)
	}
	log.Println("Order service deregistered successfully.")
}
