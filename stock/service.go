package main

import (
	"context"

	pb "example.com/oms/common/api"
	rp "example.com/oms/common/repository"
)

type StockService interface {
	AddStock(context.Context, *pb.AddStockRequest) (*pb.Stock, error)
	TakeStock(context.Context, *pb.TakeStockRequest) (*pb.Stock, error)
	GetStock(context.Context, *pb.GetStockRequest) (*pb.Stock, error)
}

type service struct {
	repository rp.Repository[*pb.Stock]
}

func NewService(repository rp.Repository[*pb.Stock]) *service {
	return &service{repository: repository}
}

func (s *service) AddStock(ctx context.Context, r *pb.AddStockRequest) (*pb.Stock, error) {
	return nil, nil
}

func (s *service) TakeStock(ctx context.Context, r *pb.TakeStockRequest) (*pb.Stock, error) {
	return nil, nil
}

func (s *service) GetStock(ctx context.Context, r *pb.GetStockRequest) (*pb.Stock, error) {
	return nil, nil
}
