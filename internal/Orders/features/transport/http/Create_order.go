package orders_http_transport

import (
	"encoding/json"
	"net/http"
)

type CreateOrderRequestDTO struct {
	Name   string `json:"order_name"`
	Price  int    `json:"order_price"`
	UserID int    `json:"user_id"`
}

func (h *OrdersHTTPHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var request CreateOrderRequestDTO

	json.NewDecoder(r.Body).Decode(&request)
}
