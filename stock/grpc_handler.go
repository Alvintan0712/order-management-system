package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "example.com/oms/common/api/protobuf"
)

type grpcHandler struct {
	pb.UnimplementedStockServiceServer

	service StockService
}

func NewGRPCHandler(grpcServer *grpc.Server, service StockService) {
	handler := &grpcHandler{service: service}
	pb.RegisterStockServiceServer(grpcServer, handler)
}

// func (h *grpcHandler) AddStock(ctx context.Context, r *pb.AddStockRequest) (*pb.Stock, error) {
// 	return h.service.AddStock(ctx, r)
// }

func (h *grpcHandler) TakeStock(ctx context.Context, r *pb.TakeStockRequest) (*pb.Stock, error) {
	return h.service.TakeStock(ctx, r)
}

func (h *grpcHandler) GetStock(ctx context.Context, r *pb.GetStockRequest) (*pb.Stock, error) {
	return h.service.GetStock(ctx, r)
}

func (h *grpcHandler) ListStocks(ctx context.Context, r *emptypb.Empty) (*pb.StockList, error) {
	stocks, err := h.service.ListStocks(ctx)
	if err != nil {
		return nil, err
	}

	stockList := &pb.StockList{
		Stocks: stocks,
	}

	return stockList, nil
}

func (h *grpcHandler) GetStocksWithMenuItem(ctx context.Context, r *emptypb.Empty) (*pb.StockMenuItemList, error) {
	stockMenuItems, err := h.service.GetStocksWithMenuItem(ctx)
	if err != nil {
		return nil, err
	}

	list := &pb.StockMenuItemList{
		Items: stockMenuItems,
	}

	return list, nil
}
