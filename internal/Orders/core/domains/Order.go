package orders_domains

import (
	pkg_errors "couriers/pkg/errors"
	"fmt"
)

var (
	UninitializedOrderID    = -1
	UnitializedOrderVersion = -1
)

type Order struct {
	ID         int
	Version    int
	Name       string
	Price      int
	IsComplete bool
	UserID     int
	CourierID  *int
}

func NewUnitializedOrder(
	name string,
	price int,
	userID int,
) Order {
	return Order{
		ID:         UninitializedOrderID,
		Version:    UnitializedOrderVersion,
		Name:       name,
		Price:      price,
		IsComplete: false,
		UserID:     userID,
		CourierID:  nil,
	}
}

func NewOrder(
	id int,
	version int,
	name string,
	price int,
	isComplete bool,
	userID int,
	courierID *int,
) Order {
	return Order{
		ID:         id,
		Version:    version,
		Name:       name,
		Price:      price,
		IsComplete: isComplete,
		UserID:     userID,
		CourierID:  courierID,
	}
}

func (o Order) Validate() error {
	if len([]rune(o.Name)) < 3 || len([]rune(o.Name)) > 20 {
		return fmt.Errorf("chars in order name must be more 3 or less 20: %w", pkg_errors.ErrInvalidArgument)
	}

	if o.Price < 0 {
		return fmt.Errorf("price must be positive or free: %w", pkg_errors.ErrInvalidArgument)
	}

	return nil
}
