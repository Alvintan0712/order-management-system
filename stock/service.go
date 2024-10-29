package main

import (
	"context"
	"errors"

	pb "example.com/oms/common/api"
	rp "example.com/oms/common/repository"
)

type StockService interface {
	AddStock(context.Context, *pb.AddStockRequest) (*pb.Stock, error)
	TakeStock(context.Context, *pb.TakeStockRequest) (*pb.Stock, error)
	GetStock(context.Context, *pb.GetStockRequest) (*pb.Stock, error)
	ListStocks(context.Context) ([]*pb.Stock, error)
}

type service struct {
	repository rp.Repository[*pb.Stock]
}

func NewService(repository rp.Repository[*pb.Stock]) *service {
	return &service{repository: repository}
}

func (s *service) AddStock(ctx context.Context, r *pb.AddStockRequest) (*pb.Stock, error) {
	stock, err := s.repository.GetById(ctx, r.ItemId)
	if err != nil {
		return nil, err
	}

	if stock != nil {
		stock.Quantity += r.Quantity
		if err := s.repository.Update(ctx, stock); err != nil {
			return nil, err
		}
	} else {
		stock = &pb.Stock{
			ItemId:   r.ItemId,
			Quantity: r.Quantity,
		}
		if err := s.repository.Create(ctx, stock); err != nil {
			return nil, err
		}
	}

	return stock, nil
}

func (s *service) TakeStock(ctx context.Context, r *pb.TakeStockRequest) (*pb.Stock, error) {
	stock, err := s.repository.GetById(ctx, r.ItemId)
	if err != nil {
		return nil, err
	}

	if stock == nil {
		return nil, errors.New("stock not enough")
	}

	if stock.Quantity < r.Quantity {
		return nil, errors.New("stock not enough")
	}

	stock.Quantity -= r.Quantity
	if err := s.repository.Update(ctx, stock); err != nil {
		return nil, err
	}

	return stock, nil
}

func (s *service) GetStock(ctx context.Context, r *pb.GetStockRequest) (*pb.Stock, error) {
	stock, err := s.repository.GetById(ctx, r.ItemId)
	if err != nil {
		return nil, err
	}

	return stock, nil
}

func (s *service) ListStocks(ctx context.Context) ([]*pb.Stock, error) {
	stocks, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}
