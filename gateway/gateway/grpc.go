package gateway

import (
	"context"
	"log"
	"time"

	pb "example.com/oms/common/api"
	"example.com/oms/common/discovery"
)

type orderGateway struct {
	registry discovery.Registry
	client   pb.OrderServiceClient
}

func NewOrderGateway(registry discovery.Registry, client pb.OrderServiceClient) *orderGateway {
	return &orderGateway{registry, client}
}

func (g *orderGateway) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.Order, error) {
	start := time.Now()
	order, err := g.client.CreateOrder(ctx, r)
	log.Printf("client create order: %v", time.Since(start))
	if err != nil {
		return nil, err
	}
	return order, nil
}
