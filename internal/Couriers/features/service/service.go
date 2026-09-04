package couriers_service

import (
	"context"
	courier_domains "couriers/internal/Couriers/core/domains"
)

type CouriersService struct {
	cache  CouriersCache
	broker CouriersBroker
}

//go:generate mockgen -source=service.go -destination=mocks/mock_couriers_repo.go -package=couriers_mocks
type CouriersCache interface {
	SaveToCache(
		ctx context.Context,
		order courier_domains.DispetchedOrder,
	) error

	GetOrder(
		ctx context.Context,
		personID int,
	) (courier_domains.DispetchedOrder, error)

	DeleteOrder(
		ctx context.Context,
		personID int,
	) error
}

type CouriersBroker interface {
	SendConfirm(
		ctx context.Context,
		orderInfo courier_domains.DispetchedOrder,
	) error
}

func NewCouriersService(
	cache CouriersCache,
	broker CouriersBroker,
) *CouriersService {
	return &CouriersService{
		cache:  cache,
		broker: broker,
	}
}
