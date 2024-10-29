package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"example.com/oms/common"
	"example.com/oms/common/discovery"
	"example.com/oms/common/discovery/consul"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	consulAddr  = common.EnvString("CONSUL_ADDR", "127.0.0.1:8500")
	serviceHost = common.EnvString("SERVICE_HOST", "127.0.0.1")
	servicePort = common.EnvString("SERVICE_PORT", "8083")
	serviceName = common.EnvString("SERVICE_NAME", "stock-service")

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
	err = registry.Register(ctx, serviceId, serviceName, serviceHost, servicePort)
	if err != nil {
		log.Fatalf("Error registering service with Consul: %v\n", err)
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

	service := NewService(repository)
	NewGRPCHandler(grpcServer, service)

	log.Printf("Stock service %s started at %s:%s\n", serviceId, serviceHost, servicePort)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err.Error())
	}
}
