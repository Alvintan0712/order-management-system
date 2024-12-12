package handler

import (
	"net/http"

	"example.com/oms/common"
	pb "example.com/oms/common/api/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CoordinatorHandler struct {
	client pb.CoordinatorServiceClient
}

func NewCoordinatorHandler(mux *http.ServeMux, conn *grpc.ClientConn) *CoordinatorHandler {
	client := pb.NewCoordinatorServiceClient(conn)
	handler := &CoordinatorHandler{client}
	handler.registerRoutes(mux)

	return handler
}

func (h *CoordinatorHandler) registerRoutes(mux *http.ServeMux) {
	// mux.HandleFunc("POST /v1/menu", h.CreateMenuItem)
}

func (h *CoordinatorHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	var req pb.CreateMenuItemWithStockRequest
	if err := common.ReadJSON(r, &req); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.client.CreateMenuItem(r.Context(), &req)
	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, response)
}
