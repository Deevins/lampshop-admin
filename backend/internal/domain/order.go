package domain

import "time"

// OrderItem описывает одну строку заказа.
type OrderItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

// OrderStatus перечисляет возможные статусы заказа.
type OrderStatus string

const (
	StatusPending    OrderStatus = "Pending"
	StatusProcessing OrderStatus = "Processing"
	StatusShipped    OrderStatus = "Shipped"
	StatusDelivered  OrderStatus = "Delivered"
)

// Order представляет модель заказа.
type Order struct {
	ID           int         `json:"id"`
	CustomerName string      `json:"customerName"`
	Items        []OrderItem `json:"items"`
	TotalPrice   float64     `json:"totalPrice"`
	Status       OrderStatus `json:"status"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}
