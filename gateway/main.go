package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/oms/common"
	pb "example.com/oms/common/api"
	"example.com/oms/common/discovery"
	"example.com/oms/common/discovery/consul"
	"example.com/oms/gateway/gateway"
	"example.com/oms/gateway/middleware"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc/grpclog"
)

var (
	httpAddr    = common.EnvString("HTTP_ADDR", ":8080")
	consulAddr  = common.EnvString("CONSUL_ADDR", "127.0.0.1:8500")
	serviceName = common.EnvString("SERVICE_NAME", "gateway")
	serviceHost = common.EnvString("SERVICE_HOST", "127.0.0.1")
	servicePort = common.EnvString("SERVICE_PORT", "8080")

	debug = common.EnvString("DEBUG", "false") == "true"
)

func main() {
	ctx := context.Background()
	if debug {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stderr))
	}

	registry, err := consul.NewRegistry(consulAddr, serviceHost, servicePort, serviceName)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	serviceId := discovery.GenerateInstanceId(serviceName)
	err = registry.Register(ctx, serviceId, serviceName, serviceHost, servicePort)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(serviceId, serviceName); err != nil {
				log.Fatalf("Failed to health check: %v\n", err)
			}
			time.Sleep(time.Second)
		}
	}()

	defer registry.Deregister(ctx, serviceId, serviceName)

	mux := http.NewServeMux()
	logHandler := middleware.Adapt(mux, middleware.Log())

	conn, err := discovery.ConnectService(ctx, "order-service", registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	if debug {
		go func() {
			for {
				log.Printf("state: %v\n", conn.GetState())
				time.Sleep(time.Second * 10)
			}
		}()
	}

	client := pb.NewOrderServiceClient(conn)
	gateway := gateway.NewOrderGateway(registry, client)

	handler := NewHandler(gateway)
	handler.registerRoutes(mux)

	log.Printf("Starting %s HTTP server at %s", serviceId, httpAddr)

	if err := http.ListenAndServe(httpAddr, logHandler); err != nil {
		log.Fatal("Failed to start http server:", err)
	}
}
