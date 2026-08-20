package dispetch_confirm_service

import (
	"context"
	dispetch_domains "couriers/internal/Dispetch/core/domains"
)

type DispetchConfirmService struct {
	db DispetchPostgresDatabase
}

type DispetchPostgresDatabase interface {
	ConfirmOrder(
		ctx context.Context,
		confirm dispetch_domains.Confirm,
	) error
}

func NewDispetchConfirmService(
	db DispetchPostgresDatabase,
) *DispetchConfirmService {
	return &DispetchConfirmService{
		db: db,
	}
}
