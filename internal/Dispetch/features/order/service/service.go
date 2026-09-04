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

//go:generate mockgen -source=service.go -destination=mocks/mock_dispetch_orders_repo.go -package=dispetch_orders_mocks
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
	cache OrdersDispetchCache,
	broker OrdersDispetchBroker,
) *OrdersDispetchService {
	return &OrdersDispetchService{
		db:     db,
		cache:  cache,
		broker: broker,
	}
}
