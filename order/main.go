package main

import (
	"context"
	"log"
	"net"

	"example.com/oms/common"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:8081")
)

func main() {
	listen, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listen.Close()

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	repository := NewRepository()
	service := NewService(repository)
	NewGRPCHandler(grpcServer, service)

	service.CreateOrder(context.Background())

	log.Println("Order server started at", grpcAddr)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err.Error())
	}
}
