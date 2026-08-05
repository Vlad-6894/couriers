package orders_repository_posgres

import orders_domains "couriers/internal/Orders/core/domains"

type OrderModel struct {
	ID         int
	Version    int
	Name       string
	Price      int
	IsComplete bool
	UserID     int
	CourierID  *int
}

type GetOrderModel struct {
	ID           int
	Version      int
	Name         string
	Price        int
	IsComplete   bool
	UserID       int
	CourierID    *int
	CourierLogin *string
}

func orderDomainFromModel(model OrderModel) orders_domains.Order {
	order := orders_domains.NewOrder(
		model.ID,
		model.Version,
		model.Name,
		model.Price,
		model.IsComplete,
		model.UserID,
		model.CourierID,
	)

	return order
}

func getOrderDomainFromModel(model GetOrderModel) orders_domains.GetOrder {
	order := orders_domains.NewGetOrder(
		model.ID,
		model.Version,
		model.Name,
		model.Price,
		model.IsComplete,
		model.UserID,
		model.CourierID,
		model.CourierLogin,
	)

	return order
}
