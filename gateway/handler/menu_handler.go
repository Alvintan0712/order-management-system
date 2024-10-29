package handler

import (
	"net/http"

	"example.com/oms/common"
	pb "example.com/oms/common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MenuHandler struct {
	client pb.MenuServiceClient
}

func NewMenuHandler(mux *http.ServeMux, conn *grpc.ClientConn) *MenuHandler {
	client := pb.NewMenuServiceClient(conn)
	handler := &MenuHandler{client}
	handler.registerRoutes(mux)

	return handler
}

func (h *MenuHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/menu", h.CreateMenuItem)
	mux.HandleFunc("GET /v1/menu/{id}", h.GetMenuItem)
	mux.HandleFunc("PUT /v1/menu/{id}", h.UpdateMenuItem)
	mux.HandleFunc("DELETE /v1/menu/{id}", h.DeleteMenuItem)
	mux.HandleFunc("GET /v1/menu", h.ListMenuItems)
}

func (h *MenuHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	var req pb.CreateMenuItemRequest
	if err := common.ReadJSON(r, &req); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	menu, err := h.client.CreateMenuItem(r.Context(), &req)
	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, menu)
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	req := pb.GetMenuItemRequest{
		Id: id,
	}

	menu, err := h.client.GetMenuItem(r.Context(), &req)
	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, menu)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req pb.UpdateMenuItemRequest
	if err := common.ReadJSON(r, &req); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if id != req.Id {
		common.WriteError(w, http.StatusBadRequest, "id not match")
		return
	}

	response, err := h.client.UpdateMenuItem(r.Context(), &req)
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

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	req := pb.DeleteMenuItemRequest{
		Id: id,
	}

	response, err := h.client.DeleteMenuItem(r.Context(), &req)
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

func (h *MenuHandler) ListMenuItems(w http.ResponseWriter, r *http.Request) {
	itemList, err := h.client.ListMenuItems(r.Context(), nil)
	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, itemList)
}
