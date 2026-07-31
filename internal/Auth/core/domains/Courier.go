package auth_domains

import (
	pkg_errors "couriers/pkg/errors"
	"fmt"
)

type Courier struct {
	ID             int
	Version        int
	Login          string
	Password       string
	City           string
	OrdersComplete int
	IsFree         bool
}

func NewCourier(
	id int,
	version int,
	login string,
	password string,
	city string,
	ordersComplete int,
	isFree bool,
) Courier {
	return Courier{
		ID:             id,
		Version:        version,
		Login:          login,
		Password:       password,
		City:           city,
		OrdersComplete: ordersComplete,
		IsFree:         isFree,
	}
}

func NewRegCourier(
	login string,
	password string,
	city string,
) Courier {
	return Courier{
		ID:             UninitializedID,
		Version:        UninitializedVersion,
		Login:          login,
		Password:       password,
		City:           city,
		OrdersComplete: 0,
		IsFree:         false,
	}
}

func (c Courier) Validate() error {
	if len([]rune(c.Login)) < 8 || len([]rune(c.Login)) > 100 {
		return fmt.Errorf(
			"login char length bigger 1 or less 8: %w",
			pkg_errors.ErrInvalidArgument,
		)
	}

	if len([]rune(c.Password)) < 8 || len([]rune(c.Password)) > 100 {
		return fmt.Errorf(
			"password char length bigger 100 or less 8: %w",
			pkg_errors.ErrInvalidArgument,
		)
	}

	if len([]rune(c.City)) < 1 || len([]rune(c.City)) > 100 {
		return fmt.Errorf(
			"login char length bigger 100 or less 1: %w",
			pkg_errors.ErrInvalidArgument,
		)
	}

	if !isUpper(c.City) {
		return fmt.Errorf("first char in the city must be upper: %w", pkg_errors.ErrInvalidArgument)
	}

	return nil
}
