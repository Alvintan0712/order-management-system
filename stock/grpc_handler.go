package main

import (
	"context"

	"google.golang.org/grpc"

	pb "example.com/oms/common/api"
)

type grpcHandler struct {
	pb.UnimplementedStockServiceServer

	service StockService
}

func NewGRPCHandler(grpcServer *grpc.Server, service StockService) {
	handler := &grpcHandler{service: service}
	pb.RegisterStockServiceServer(grpcServer, handler)
}

func (h *grpcHandler) AddStock(ctx context.Context, r *pb.AddStockRequest) (*pb.Stock, error) {
	return h.service.AddStock(ctx, r)
}

func (h *grpcHandler) TakeStock(ctx context.Context, r *pb.TakeStockRequest) (*pb.Stock, error) {
	return h.service.TakeStock(ctx, r)
}

func (h *grpcHandler) GetStock(ctx context.Context, r *pb.GetStockRequest) (*pb.Stock, error) {
	return h.service.GetStock(ctx, r)
}
