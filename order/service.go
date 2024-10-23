package main

import (
	"context"
	"log"

	pb "example.com/oms/common/api"
)

type service struct {
	repository OrderRepository
}

func NewService(repository OrderRepository) *service {
	return &service{repository}
}

func (s *service) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Println("New order received!")
	order := &pb.Order{
		Id: "42",
	}
	return order, nil
}
