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
	resp, err := h.service.CreateMenuItemWithStock(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
