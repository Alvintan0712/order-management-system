package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"example.com/oms/common"
	pb "example.com/oms/common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StockHandler struct {
	client pb.StockServiceClient
}

func NewStockHandler(mux *http.ServeMux, conn *grpc.ClientConn) *StockHandler {
	client := pb.NewStockServiceClient(conn)
	handler := &StockHandler{client}
	handler.registerRoutes(mux)

	return handler
}

func (h *StockHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/stock", h.AddStock)
	mux.HandleFunc("GET /v1/stock/{id}", h.GetStock)
	mux.HandleFunc("PUT /v1/stock/{id}", h.TakeStock)
	mux.HandleFunc("GET /v1/stock", h.ListStocks)
	mux.HandleFunc("GET /v1/stock/menu", h.ListStocksWithMenu)
}

func (h *StockHandler) AddStock(w http.ResponseWriter, r *http.Request) {
	var req pb.AddStockRequest
	if err := common.ReadJSON(r, &req); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	menu, err := h.client.AddStock(r.Context(), &req)
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

func (h *StockHandler) GetStock(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	req := pb.GetStockRequest{
		ItemId: id,
	}

	menu, err := h.client.GetStock(r.Context(), &req)
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

func (h *StockHandler) TakeStock(w http.ResponseWriter, r *http.Request) {
	var req pb.TakeStockRequest
	if err := common.ReadJSON(r, &req); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.client.TakeStock(r.Context(), &req)
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

func (h *StockHandler) ListStocks(w http.ResponseWriter, r *http.Request) {
	stockList, err := h.client.ListStocks(r.Context(), nil)
	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, stockList)
}

func (h *StockHandler) ListStocksWithMenu(w http.ResponseWriter, r *http.Request) {
	list, err := h.client.GetStocksWithMenuItem(r.Context(), nil)

	jsonData, _ := json.Marshal(list)
	log.Println(string(jsonData))

	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, list)
}
