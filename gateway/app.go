package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/oms/common/discovery"
	"example.com/oms/common/discovery/consul"
	"example.com/oms/gateway/middleware"
	"example.com/oms/gateway/service"
	"google.golang.org/grpc/grpclog"
)

type Config struct {
	httpAddr   string
	consulAddr string
	Id         string
	Name       string
	Host       string
	Port       string
}

type app struct {
	context      context.Context
	config       *Config
	registry     *consul.Registry
	mux          *http.ServeMux
	handler      http.Handler
	grpcServices []service.GRPCService
}

func NewApp(config *Config) *app {
	config.Id = discovery.GenerateInstanceId(config.Name)
	app := &app{
		context: context.Background(),
		config:  config,
		mux:     http.NewServeMux(),
	}

	if debug {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stderr))
	}

	app.registerServer()
	app.setupMiddleware()
	app.setupServices()

	return app
}

func (a *app) Listen() error {
	log.Printf("Starting %s HTTP server at %s", a.config.Id, a.config.httpAddr)
	if err := http.ListenAndServe(a.config.httpAddr, a.handler); err != nil {
		return err
	}
	return nil
}

func (a *app) Close() {
	a.registry.Deregister(a.context, a.config.Id, a.config.Name)
	for _, service := range a.grpcServices {
		service.Close()
	}
}

func (a *app) registerServer() {
	registry, err := consul.NewRegistry(a.config.consulAddr, a.config.Host, a.config.Port, a.config.Name)
	if err != nil {
		log.Fatalf("error in create new client: %v\n", err)
	}

	err = registry.Register(a.context, a.config.Id, a.config.Name, a.config.Host, a.config.Port)
	if err != nil {
		log.Fatalf("error in register server: %v\n", err)
	}

	// health check
	go func() {
		for {
			if err := registry.HealthCheck(a.config.Id, a.config.Name); err != nil {
				log.Fatalf("Failed to health check: %v\n", err)
			}
			time.Sleep(time.Second)
		}
	}()

	a.registry = registry
}

func (a *app) setupMiddleware() {
	handler := middleware.Adapt(a.mux, middleware.Log())
	a.handler = handler
}

func (a *app) setupServices() {
	coordinatorService, err := service.NewCoordinatorService(a.context, a.mux, "coordinator", a.registry)
	if err != nil {
		log.Fatalf("coordinator service create failed: %v\n", err)
	}
	a.grpcServices = append(a.grpcServices, coordinatorService.GRPCService)

	orderService, err := service.NewOrderService(a.context, a.mux, "order-service", a.registry)
	if err != nil {
		log.Fatalf("order service create failed: %v\n", err)
	}
	a.grpcServices = append(a.grpcServices, orderService.GRPCService)

	menuService, err := service.NewMenuService(a.context, a.mux, "menu-service", a.registry)
	if err != nil {
		log.Fatalf("menu service create failed: %v\n", err)
	}
	a.grpcServices = append(a.grpcServices, menuService.GRPCService)

	stockService, err := service.NewStockService(a.context, a.mux, "stock-service", a.registry)
	if err != nil {
		log.Fatalf("stock service create failed: %v\n", err)
	}
	a.grpcServices = append(a.grpcServices, stockService.GRPCService)
}
