package orders_http_transport

type OrdersHTTPHandler struct {
	ordersService OrdersService
}

type OrdersService interface{}

func NewOrdersHTTPHandler(ordersService OrdersService) *OrdersHTTPHandler {
	return &OrdersHTTPHandler{
		ordersService: ordersService,
	}
}
