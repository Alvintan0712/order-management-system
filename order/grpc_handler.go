package main

import (
	"context"
	"log"

	pb "example.com/oms/common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer

	service OrderService
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrderService) {
	handler := &grpcHandler{}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received! Order %v\n", r)
	order := &pb.Order{
		Id: "42",
	}
	return order, nil
}
