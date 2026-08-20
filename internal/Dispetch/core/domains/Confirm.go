package dispetch_domains

type Confirm struct {
	OrderID   int
	CourierID int
}

func NewConfirm(
	orderID int,
	courierID int,
) Confirm {
	return Confirm{
		OrderID:   orderID,
		CourierID: courierID,
	}
}
