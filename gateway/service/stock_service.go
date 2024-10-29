package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"example.com/oms/common/discovery"
	"example.com/oms/common/discovery/consul"
	"example.com/oms/gateway/handler"
)

type StockService struct {
	GRPCService
	handler *handler.StockHandler
}

func NewStockService(ctx context.Context, mux *http.ServeMux, name string, gatewayRegistry *consul.Registry) (*StockService, error) {
	conn, err := discovery.ConnectService(ctx, name, gatewayRegistry)
	if err != nil {
		return nil, err
	}

	handler := handler.NewStockHandler(mux, conn)

	if debug {
		go func() {
			for {
				log.Printf("stock service state: %v\n", conn.GetState())
				time.Sleep(time.Second * 10)
			}
		}()
	}

	service := &StockService{
		GRPCService: GRPCService{
			Service: Service{
				Name: name,
			},
			Registry:   gatewayRegistry,
			Connection: conn,
		},
		handler: handler,
	}

	return service, nil
}
