package main

import (
	"context"
	"log"

	pb "example.com/oms/common/api"
)

type MenuService interface {
	CreateMenuItem(context.Context, *pb.CreateMenuItemRequest) (*pb.MenuItem, error)
}

type service struct {
	repository MenuRepository
}

func NewService(repository MenuRepository) *service {
	return &service{repository: repository}
}

func (s *service) CreateMenuItem(ctx context.Context, r *pb.CreateMenuItemRequest) (*pb.MenuItem, error) {
	log.Println("Create menu item")
	menu := &pb.MenuItem{
		Id:        "123",
		Name:      "burger",
		UnitPrice: 599,
		Currency:  "MYR",
	}
	return menu, nil
}
