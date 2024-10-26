package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"example.com/oms/common"
	"example.com/oms/common/discovery"
	"example.com/oms/common/discovery/consul"
	"example.com/oms/gateway/gateway"
	"example.com/oms/gateway/middleware"
	_ "github.com/joho/godotenv/autoload"
)

var (
	httpAddr    = common.EnvString("HTTP_ADDR", ":8080")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	serviceName = common.EnvString("SERVICE_NAME", "gateway")
	serviceHost = common.EnvString("SERVICE_HOST", "localhost")
	servicePort = common.EnvString("SERVICE_PORT", "8080")
)

func main() {
	ctx := context.Background()

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

	gateway := gateway.NewOrderGateway(ctx, registry)

	handler := NewHandler(gateway)
	handler.registerRoutes(mux)

	log.Printf("Starting %s HTTP server at %s", serviceId, httpAddr)

	if err := http.ListenAndServe(httpAddr, logHandler); err != nil {
		log.Fatal("Failed to start http server:", err)
	}
}
