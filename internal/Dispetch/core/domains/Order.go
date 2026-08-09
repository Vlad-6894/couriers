package dispetch_domains

type Order struct {
	ID         int
	Version    int
	Name       string
	Price      int
	IsComplete bool
	City       string
	UserID     int
}

func NewOrder(
	id int,
	version int,
	name string,
	price int,
	isComplete bool,
	city string,
	userID int,
) Order {
	return Order{
		ID:         id,
		Version:    version,
		Name:       name,
		Price:      price,
		IsComplete: isComplete,
		City:       city,
		UserID:     userID,
	}
}
