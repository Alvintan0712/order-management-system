package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"example.com/oms/common"
	pb "example.com/oms/common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	client pb.OrderServiceClient
}

func registerRoutes(mux *http.ServeMux, handler *OrderHandler) {
	mux.HandleFunc("POST /v1/customer/{customerId}/order", handler.CreateOrder)
}

func NewOrderHandler(mux *http.ServeMux, conn *grpc.ClientConn) *OrderHandler {
	client := pb.NewOrderServiceClient(conn)
	handler := &OrderHandler{client}
	registerRoutes(mux, handler)

	return handler
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
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
	order, err := h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
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
