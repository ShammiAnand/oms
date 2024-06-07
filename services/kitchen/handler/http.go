package kitchen

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/shammianand/oms/services/common/genproto/orders"
	"github.com/shammianand/oms/services/common/util"
	"google.golang.org/grpc"
)

type Handler struct {
	conn *grpc.ClientConn
}

func NewHandler(conn *grpc.ClientConn) *Handler {
	return &Handler{conn: conn}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /orders/{CustomerID}", h.getOrdersByCustomerID)
	router.HandleFunc("POST /orders", h.createOrder)
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	c := orders.NewOrderServiceClient(h.conn)
	var createOrderReq orders.CreateOrderRequest
	err := util.ParseJSON(r, &createOrderReq)
	if err != nil {
		util.WriteError(
			w, http.StatusFailedDependency, err,
		)
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	createOrderRes, err := c.CreateOrder(ctx, &createOrderReq)
	if err != nil {
		util.WriteError(
			w, http.StatusFailedDependency, err,
		)
	}
	util.WriteJSON(w, http.StatusOK, createOrderRes)
}

func (h *Handler) getOrdersByCustomerID(w http.ResponseWriter, r *http.Request) {
	custID, err := strconv.Atoi(r.PathValue("CustomerID"))
	if err != nil {
		util.WriteError(w, http.StatusFailedDependency, err)
	}
	c := orders.NewOrderServiceClient(h.conn)

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	res, err := c.GetOrders(ctx, &orders.GetOrdersRequest{
		CustomerID: int32(custID),
	})
	if err != nil {
		log.Fatalf("client error: %v", err)
	}
	util.WriteJSON(w, http.StatusOK, res)

}
