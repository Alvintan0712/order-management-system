package gateway

import (
	"context"

	pb "example.com/oms/common/api"
)

type OrderGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
}
