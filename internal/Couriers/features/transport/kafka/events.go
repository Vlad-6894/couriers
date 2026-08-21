package couriers_kafka_transport

type DispetchedOrderEventForCourier struct {
	OrderID   int `json:"Order_id"`
	Version   int `json:"Version"`
	CourierID int `json:"Courier_id"`
}
