package service

import (
	"context"
	"log"

	"github.com/shammianand/oms/services/common/genproto/orders"
)

var ordersDb = make([]*orders.Order, 0)

type OrderService struct {
	// store
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *orders.Order) error {
	ordersDb = append(ordersDb, order)
	return nil
}

func (s *OrderService) GetOrders(ctx context.Context, customerID int32) []*orders.Order {
	log.Println("this method is getting called", customerID)
	return ordersDb
}
