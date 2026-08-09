package dispetch_kafka_transport

type OrderCreatedEvent struct {
	ID         int
	Version    int
	Name       string
	Price      int
	IsComplete bool
	City       string
	UserID     int
}
