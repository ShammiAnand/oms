package handler

import (
	"context"
	"math/rand"

	"github.com/shammianand/oms/services/common/genproto/orders"
	"github.com/shammianand/oms/services/orders/types"
	"google.golang.org/grpc"
)

type OrdersGrpcHandler struct {
	ordersService types.OrderService
	orders.UnimplementedOrderServiceServer
}

func NewGrpcOrdersService(grpc *grpc.Server, ordersService types.OrderService) {
	gRPCHandler := &OrdersGrpcHandler{
		ordersService: ordersService,
	}

	// register the OrderServiceServer
	orders.RegisterOrderServiceServer(grpc, gRPCHandler)
}

func (h *OrdersGrpcHandler) GetOrders(ctx context.Context, req *orders.GetOrdersRequest) (*orders.GetOrderResponse, error) {
	o := h.ordersService.GetOrders(ctx, req.CustomerID)
	res := &orders.GetOrderResponse{
		Orders: o,
	}

	return res, nil
}

func (h *OrdersGrpcHandler) CreateOrder(ctx context.Context, req *orders.CreateOrderRequest) (*orders.CreateOrderResponse, error) {

	order := &orders.Order{
		OrderID:    int32(rand.Intn(10000)),
		CustomerID: req.CustomerID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
	}

	err := h.ordersService.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	res := &orders.CreateOrderResponse{
		Status: "success",
	}

	return res, nil
}
