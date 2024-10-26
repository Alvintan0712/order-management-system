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

func NewOrderGateway(ctx context.Context, registry discovery.Registry) *orderGateway {
	conn, err := discovery.ConnectService(ctx, "order-service", registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	client := pb.NewOrderServiceClient(conn)

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
	// return client.CreateOrder(ctx, r)
}
