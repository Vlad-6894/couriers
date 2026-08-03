package orders_domains

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
	UserID     int
	CourierID  *int
}

func NewUnitializedOrder(
	name string,
	price int,
	userID int,
) Order {
	return Order{
		ID:         UninitializedOrderID,
		Version:    UnitializedOrderVersion,
		Name:       name,
		Price:      price,
		IsComplete: false,
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
	userID int,
	courierID *int,
) Order {
	return Order{
		ID:         id,
		Version:    version,
		Name:       name,
		Price:      price,
		IsComplete: isComplete,
		UserID:     userID,
		CourierID:  courierID,
	}
}
