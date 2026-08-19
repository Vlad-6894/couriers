package dispetch_kafka_repository

type DispetchedOrderEvent struct {
	OrderID   int `json:"Order_id"`
	Version   int `json:"Version"`
	CourierID int `json:"Courier_id"`
}
