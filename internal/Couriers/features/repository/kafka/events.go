package couriers_repository_kafka

type ConfirmEvent struct {
	OrderID   int `json:"order_id"`
	CourierID int `json:"courier_id"`
}
