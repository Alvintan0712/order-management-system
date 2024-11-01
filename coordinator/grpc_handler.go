package main

import (
	"context"

	pb "example.com/oms/common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedCoordinatorServiceServer

	service CoordinatorService
}

func NewGRPCHandler(grpcServer *grpc.Server, service CoordinatorService) {
	handler := &grpcHandler{service: service}
	pb.RegisterCoordinatorServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateMenuItem(ctx context.Context, r *pb.CreateMenuItemWithStockRequest) (*pb.CreateMenuItemResponse, error) {
	response := &pb.CreateMenuItemResponse{
		Success: true,
		Message: "Create menu item success",
	}

	return response, nil
}
