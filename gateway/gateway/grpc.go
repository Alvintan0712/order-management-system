package gateway

import (
	"context"
	"log"

	pb "example.com/oms/common/api"
	"example.com/oms/common/discovery"
)

type orderGateway struct {
	registry discovery.Registry
}

func NewOrderGateway(registry discovery.Registry) *orderGateway {
	return &orderGateway{registry}
}

func (g *orderGateway) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.Order, error) {
	conn, err := discovery.ConnectService(ctx, "order-service", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	client := pb.NewOrderServiceClient(conn)

	return client.CreateOrder(ctx, r)
}
