package types

import (
	"context"

	"github.com/shammianand/oms/services/common/genproto/orders"
)

type OrderService interface {
	CreateOrder(context.Context, *orders.Order) error
	GetOrders(context.Context, int32) []*orders.Order
}
