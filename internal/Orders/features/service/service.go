package orders_service

import (
	"context"
	orders_domains "couriers/internal/Orders/core/domains"
)

type OrdersService struct {
	db     DatabasePostgres
	broker BrokerKafka
}

type DatabasePostgres interface {
	CreateOrder(
		ctx context.Context,
		order orders_domains.Order,
	) (orders_domains.Order, error)

	GetOrders(
		ctx context.Context,
		personID int,
		limit *int,
		offset *int,
	) ([]orders_domains.GetOrder, error)
}

type BrokerKafka interface {
	SendOrder(
		ctx context.Context,
		order orders_domains.Order,
	) error
}

func NewOrdersService(
	db DatabasePostgres,
	broker BrokerKafka,
) *OrdersService {
	return &OrdersService{
		db:     db,
		broker: broker,
	}
}
