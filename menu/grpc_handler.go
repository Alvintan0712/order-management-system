package main

import (
	"context"
	"log"

	pb "example.com/oms/common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedMenuServiceServer

	service MenuService
}

func NewGRPCHandler(grpcServer *grpc.Server, service MenuService) {
	handler := &grpcHandler{service: service}
	pb.RegisterMenuServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateMenuItem(ctx context.Context, r *pb.CreateMenuItemRequest) (*pb.MenuItem, error) {
	log.Printf("Create menu item, req: %v\n", r)
	menu, err := h.service.CreateMenuItem(ctx, r)
	if err != nil {
		return nil, err
	}

	return menu, nil
}
