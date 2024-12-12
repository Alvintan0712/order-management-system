package service

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "example.com/oms/common/api/protobuf"
	"example.com/oms/common/discovery"
	"example.com/oms/common/discovery/consul"
	"example.com/oms/gateway/handler"
)

type OrderService struct {
	GRPCService
	client  *pb.OrderServiceClient
	handler *handler.OrderHandler
}

func NewOrderService(ctx context.Context, mux *http.ServeMux, name string, gatewayRegistry *consul.Registry) (*OrderService, error) {
	conn, err := discovery.ConnectService(ctx, name, gatewayRegistry)
	if err != nil {
		return nil, err
	}

	client := pb.NewOrderServiceClient(conn)
	handler := handler.NewOrderHandler(mux, conn)

	if debug {
		go func() {
			for {
				log.Printf("order service state: %v\n", conn.GetState())
				time.Sleep(time.Second * 10)
			}
		}()
	}

	service := &OrderService{
		GRPCService: GRPCService{
			Service: Service{
				Name: name,
			},
			Registry:   gatewayRegistry,
			Connection: conn,
		},
		client:  &client,
		handler: handler,
	}

	return service, nil
}
