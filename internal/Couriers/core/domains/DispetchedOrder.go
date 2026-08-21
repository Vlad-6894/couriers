package courier_domains

type DispetchedOrder struct {
	OrderID   int
	Version   int
	CourierID int
}

func NewDispetchedOrder(
	orderID int,
	version int,
	courierID int,
) DispetchedOrder {
	return DispetchedOrder{
		OrderID:   orderID,
		Version:   version,
		CourierID: courierID,
	}
}
