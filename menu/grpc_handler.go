package main

import (
	"context"
	"log"

	pb "example.com/oms/common/api/protobuf"
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
	log.Println("menu handler: create menu item")
	menu, err := h.service.CreateMenuItem(ctx, r)
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (h *grpcHandler) GetMenuItem(ctx context.Context, r *pb.GetMenuItemRequest) (*pb.MenuItem, error) {
	menu, err := h.service.GetMenuItem(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (h *grpcHandler) UpdateMenuItem(ctx context.Context, r *pb.UpdateMenuItemRequest) (*pb.UpdateMenuItemResponse, error) {
	err := h.service.UpdateMenuItem(ctx, r)
	if err != nil {
		return nil, err
	}

	response := &pb.UpdateMenuItemResponse{
		Message: "menu item updated successfully!",
	}

	return response, nil
}

func (h *grpcHandler) DeleteMenuItem(ctx context.Context, r *pb.DeleteMenuItemRequest) (*pb.DeleteMenuItemResponse, error) {
	err := h.service.DeleteMenuItem(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	response := &pb.DeleteMenuItemResponse{
		Message: "menu item deleted successfully!",
	}

	return response, nil
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
