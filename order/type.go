package main

import (
	"context"

	pb "example.com/oms/common/api"
)

type OrderService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) error
}

type OrderRepository interface {
	Create(context.Context) error
}
