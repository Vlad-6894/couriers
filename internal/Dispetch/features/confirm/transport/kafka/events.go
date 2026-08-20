package dispetch_kafka_confirm_transport

type ConfirmEvent struct {
	OrderID   int `json:"order_id"`
	CourierID int `json:"courier_id"`
}
