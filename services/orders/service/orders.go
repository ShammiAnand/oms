package service

import (
	"context"
	"log"

	"github.com/shammianand/oms/services/common/genproto/orders"
)

var ordersDb = make([]*orders.Order, 0)

type OrderService struct {
	// TODO: add mongo store
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *orders.Order) error {
	ordersDb = append(ordersDb, order)
	return nil
}

func (s *OrderService) GetOrders(ctx context.Context, customerID int32) []*orders.Order {
	log.Println(customerID)
	customerOrders := make([]*orders.Order, 0)
	for _, item := range ordersDb {
		if item.CustomerID == customerID {
			log.Println("matched")
			customerOrders = append(customerOrders, item)
		}
	}
	return customerOrders
}
