package dispetch_service

import (
	"context"
	dispetch_domains "couriers/internal/Dispetch/core/domains"
)

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

	DoBusy(
		ctx context.Context,
		courierID int,
		version int,
	) error

	GetFreeCouriers(
		ctx context.Context,
	) ([]dispetch_domains.CourierInfo, error)
}

type OrdersDispetchBroker interface {
	SendDispetchedOrder(
		ctx context.Context,
		orderID int,
		version int,
		courierID int,
		city string,
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

	UpdateCache(
		ctx context.Context,
		couriersByGroupCity map[string]map[int]dispetch_domains.FreeCourierInfo,
	) error
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
