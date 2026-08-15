package dispetch_service

import "context"

type OrdersDispetchService struct {
	db     OrdersDispetchDatabase
	cache  OrdersDispetchCache
	broker OrdersDispetchBroker
}

type OrdersDispetchDatabase interface {
	SearchCourier(
		ctx context.Context,
		city string,
	) (int, int, error)
}

type OrdersDispetchBroker interface {
	SendDispetchedOrder(
		ctx context.Context,
		orderID int,
		version int,
		courierID int,
	) error
}

type OrdersDispetchCache interface {
	CheckUnique(
		ctx context.Context,
		orderId int,
	) (bool, error)

	SearchCourier(
		ctx context.Context,
		city string,
	) (int, int, error)
}

func NewOrdersDispetchService(
	db OrdersDispetchDatabase,
	broker OrdersDispetchBroker,
) *OrdersDispetchService {
	return &OrdersDispetchService{
		db:     db,
		broker: broker,
	}
}
