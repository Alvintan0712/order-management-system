package main

import (
	"context"

	pb "example.com/oms/common/api"
	"example.com/oms/common/discovery"
	"example.com/oms/common/discovery/consul"
)

type CoordinatorService interface {
	CreateMenuItemWithStock(context.Context, *pb.CreateMenuItemWithStockRequest) (*pb.CreateMenuItemResponse, error)
}

type service struct {
	menuClient  pb.MenuServiceClient
	stockClient pb.StockServiceClient
}

func NewService(ctx context.Context, registry *consul.Registry) (*service, error) {
	conn, err := discovery.ConnectService(ctx, "menu-service", registry)
	if err != nil {
		return nil, err
	}
	menuClient := pb.NewMenuServiceClient(conn)

	conn, err = discovery.ConnectService(ctx, "stock-service", registry)
	if err != nil {
		return nil, err
	}
	stockClient := pb.NewStockServiceClient(conn)

	return &service{menuClient: menuClient, stockClient: stockClient}, nil
}

func (s *service) CreateMenuItemWithStock(ctx context.Context, r *pb.CreateMenuItemWithStockRequest) (*pb.CreateMenuItemResponse, error) {
	menuItem, err := s.menuClient.CreateMenuItem(ctx, &pb.CreateMenuItemRequest{
		Name:      r.Name,
		UnitPrice: r.UnitPrice,
		Currency:  r.Currency,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.stockClient.AddStock(ctx, &pb.AddStockRequest{
		ItemId:   menuItem.Id,
		Quantity: 0,
	})
	if err != nil {
		_, deleteErr := s.menuClient.DeleteMenuItem(ctx, &pb.DeleteMenuItemRequest{
			Id: menuItem.Id,
		})
		if deleteErr != nil {
			return nil, deleteErr
		}
		return nil, err
	}

	response := &pb.CreateMenuItemResponse{
		Success: true,
		Message: "Create menu item successfully",
	}

	return response, nil
}
