package main

import (
	"context"
	"errors"

	pb "example.com/oms/common/api/protobuf"
	"example.com/oms/common/discovery"
	"example.com/oms/common/discovery/consul"
	rp "example.com/oms/common/repository"
	"google.golang.org/protobuf/proto"
)

type StockService interface {
	AddStock(context.Context, *pb.AddStockRequest) (*pb.Stock, error)
	TakeStock(context.Context, *pb.TakeStockRequest) (*pb.Stock, error)
	GetStock(context.Context, *pb.GetStockRequest) (*pb.Stock, error)
	ListStocks(context.Context) ([]*pb.Stock, error)
	GetStocksWithMenuItem(context.Context) ([]*pb.StockMenuItem, error)
}

type service struct {
	repository rp.Repository[*pb.Stock]

	menuClient pb.MenuServiceClient
}

func NewService(ctx context.Context, repository rp.Repository[*pb.Stock], registry *consul.Registry) (*service, error) {
	conn, err := discovery.ConnectService(ctx, "menu-service", registry)
	if err != nil {
		return nil, err
	}

	menuClient := pb.NewMenuServiceClient(conn)

	return &service{repository: repository, menuClient: menuClient}, nil
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

func (s *service) GetStocksWithMenuItem(ctx context.Context) ([]*pb.StockMenuItem, error) {
	menuList, err := s.menuClient.ListMenuItems(ctx, nil)
	if err != nil {
		return nil, err
	}

	stocks, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}
	itemQuantityMap := make(map[string]int32)
	for _, stock := range stocks {
		itemQuantityMap[stock.ItemId] = stock.Quantity
	}

	stockMenuItems := make([]*pb.StockMenuItem, len(menuList.Items))
	for i, menuItem := range menuList.Items {
		quantity, ok := itemQuantityMap[menuItem.Id]
		if !ok {
			quantity = 0
		}

		stockMenuItems[i] = &pb.StockMenuItem{
			ItemId:   menuItem.Id,
			Quantity: proto.Int32(quantity),
			Item:     menuItem,
		}
	}

	return stockMenuItems, nil
}
