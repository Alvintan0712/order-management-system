package handler

import (
	"net/http"

	pb "example.com/oms/common/api/protobuf"
	"google.golang.org/grpc"
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

}
