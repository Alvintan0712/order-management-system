package main

import (
	"log"
	"net/http"

	"example.com/oms/common"
	pb "example.com/oms/common/api"
)

type handler struct {
	orderClient pb.OrderServiceClient
}

func NewHandler(orderClient pb.OrderServiceClient) *handler {
	return &handler{orderClient}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customer/{customerId}/order", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customerId")
	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Println("Forward to order service")
	_, err := h.orderClient.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerId: customerId,
		Items:      items,
	})

	if err != nil {
		log.Fatal("Something wrong:", err)
	}
}
