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

type CoordinatorService struct {
	GRPCService
	handler *handler.CoordinatorHandler
}

func NewCoordinatorService(ctx context.Context, mux *http.ServeMux, name string, gatewayRegistry *consul.Registry) (*CoordinatorService, error) {
	conn, err := discovery.ConnectService(ctx, name, gatewayRegistry)
	if err != nil {
		return nil, err
	}

	handler := handler.NewCoordinatorHandler(mux, conn)

	if debug {
		go func() {
			for {
				log.Printf("coordinator service state: %v\n", conn.GetState())
				time.Sleep(time.Second * 10)
			}
		}()
	}

	service := &CoordinatorService{
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
