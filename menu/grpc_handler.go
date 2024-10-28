package main

import (
	"log"

	"google.golang.org/grpc"
)

type grpcHandler struct {
	service MenuService
}

func NewGRPCHandler(grpcServer *grpc.Server, service MenuService) {
	handler := &grpcHandler{service: service}
	log.Println(handler)
	// register handler
}
