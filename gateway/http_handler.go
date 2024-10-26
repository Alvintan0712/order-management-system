package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"example.com/oms/common"
	pb "example.com/oms/common/api"
	"example.com/oms/gateway/gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	orderGateway gateway.OrderGateway
}

func NewHandler(orderGateway gateway.OrderGateway) *handler {
	return &handler{orderGateway}
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

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	start := time.Now()
	order, err := h.orderGateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerId: customerId,
		Items:      items,
	})
	log.Printf("create order: %v\n", time.Since(start))

	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, order)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrNoItems
	}

	for _, i := range items {
		if i.Id == "" {
			return errors.New("item Id is required")
		}

		if i.Quantity <= 0 {
			return errors.New("item must have a valid quantity")
		}
	}

	return nil
}
