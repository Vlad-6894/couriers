package dispetch_kafka_transport

type OrderCreatedEvent struct {
	ID         int    `json:"order_id"`
	Version    int    `json:"order_version"`
	Name       string `json:"order_name"`
	Price      int    `json:"order_price"`
	IsComplete bool   `json:"is_complete"`
	City       string `json:"city"`
	UserID     int    `json:"user_id"`
}
