package orders_domains

import (
	pkg_errors "couriers/pkg/errors"
	"fmt"
	"unicode"
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
	City       string
	UserID     int
	CourierID  *int
}

type GetOrder struct {
	ID           int
	Version      int
	Name         string
	Price        int
	IsComplete   bool
	City         string
	UserLogin    string
	UserID       int
	CourierID    *int
	CourierLogin *string
}

func NewUnitializedOrder(
	name string,
	price int,
	city string,
	userID int,
) Order {
	return Order{
		ID:         UninitializedOrderID,
		Version:    UnitializedOrderVersion,
		Name:       name,
		Price:      price,
		IsComplete: false,
		City:       city,
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
	city string,
	userID int,
	courierID *int,
) Order {
	return Order{
		ID:         id,
		Version:    version,
		Name:       name,
		Price:      price,
		IsComplete: isComplete,
		City:       city,
		UserID:     userID,
		CourierID:  courierID,
	}
}

func NewGetOrder(
	id int,
	version int,
	name string,
	price int,
	isComplete bool,
	city string,
	userLogin string,
	userID int,
	courierID *int,
	courierLogin *string,
) GetOrder {
	return GetOrder{
		ID:           id,
		Version:      version,
		Name:         name,
		Price:        price,
		IsComplete:   isComplete,
		City:         city,
		UserLogin:    userLogin,
		UserID:       userID,
		CourierID:    courierID,
		CourierLogin: courierLogin,
	}
}

func (o Order) Validate() error {
	if len([]rune(o.Name)) < 3 || len([]rune(o.Name)) > 20 {
		return fmt.Errorf("chars in order name must be more 3 or less 20: %w", pkg_errors.ErrInvalidArgument)
	}

	if o.Price < 0 {
		return fmt.Errorf("price must be positive or free: %w", pkg_errors.ErrInvalidArgument)
	}

	if len([]rune(o.City)) < 1 || len([]rune(o.City)) > 100 {
		return fmt.Errorf(
			"login char length bigger 100 or less 1: %w",
			pkg_errors.ErrInvalidArgument,
		)
	}

	if !isUpper(o.City) {
		return fmt.Errorf("first char in the city must be upper: %w", pkg_errors.ErrInvalidArgument)
	}

	return nil
}

func isUpper(str string) bool {
	r := []rune(str)[0]
	return unicode.IsUpper(r)
}
