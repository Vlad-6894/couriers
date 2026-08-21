package couriers_transport

import (
	"context"
	courier_domains "couriers/internal/Couriers/core/domains"
)

type CouriersService interface {
	SaveToCache(
		ctx context.Context,
		order courier_domains.DispetchedOrder,
	) error

	GetOrder(
		ctx context.Context,
		personID int,
	) (courier_domains.DispetchedOrder, error)

	Confirm(
		ctx context.Context,
		personID int,
	) error
}
