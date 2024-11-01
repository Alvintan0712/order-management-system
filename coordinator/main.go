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
	"google.golang.org/grpc"
)

var (
	consulAddr  = common.EnvString("CONSUL_ADDR", "127.0.0.1:8500")
	serviceHost = common.EnvString("SERVICE_HOST", "127.0.0.1")
	servicePort = common.EnvString("SERVICE_PORT", "8800")
	serviceName = common.EnvString("SERVICE_NAME", "coordinator")

	debug = common.EnvString("DEBUG", "false") == "true"
)

func main() {
	ctx := context.Background()

	if debug {
		//
	}

	registry, err := consul.NewRegistry(consulAddr, serviceHost, servicePort, serviceName)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v\n", err)
	}

	serviceId := discovery.GenerateInstanceId(serviceName)
	err = registry.Register(ctx, serviceId, serviceName, serviceHost, servicePort)
	if err != nil {
		log.Fatalf("Error registering service with Consul: %v\n", err)
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

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", serviceHost, servicePort))
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	defer listen.Close()

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	service := NewService()
	NewGRPCHandler(grpcServer, service)

	log.Printf("Coordinator %s started at %s:%s\n", serviceId, serviceHost, servicePort)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err.Error())
	}
}
