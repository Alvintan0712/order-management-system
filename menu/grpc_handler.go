package main

import (
	"context"

	pb "example.com/oms/common/api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
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
	menu, err := h.service.CreateMenuItem(ctx, r)
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (h *grpcHandler) ListMenuItems(ctx context.Context, r *emptypb.Empty) (*pb.MenuItemList, error) {
	items, err := h.service.ListMenuItems(ctx)
	if err != nil {
		return nil, err
	}

	itemList := &pb.MenuItemList{
		Items: items,
	}
	return itemList, nil
}
