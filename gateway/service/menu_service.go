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

type MenuService struct {
	GRPCService
	handler *handler.MenuHandler
}

func NewMenuService(ctx context.Context, mux *http.ServeMux, name string, gatewayRegistry *consul.Registry) (*MenuService, error) {
	conn, err := discovery.ConnectService(ctx, name, gatewayRegistry)
	if err != nil {
		return nil, err
	}

	handler := handler.NewMenuHandler(mux, conn)

	if debug {
		go func() {
			for {
				log.Printf("menu service state: %v\n", conn.GetState())
				time.Sleep(time.Second * 10)
			}
		}()
	}

	service := &MenuService{
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
