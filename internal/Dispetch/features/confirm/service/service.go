package dispetch_confirm_service

import (
	"context"
	dispetch_domains "couriers/internal/Dispetch/core/domains"
)

type DispetchConfirmService struct {
	db DispetchPostgresDatabase
}

//go:generate mockgen -source=service.go -destination=mocks/mock_dispetch_confirm_repo.go -package=dispetch_confirm_mocks
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
