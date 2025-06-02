package domain

import "time"

// Product представляет модель товара, включая JSONB-атрибуты.
type Product struct {
	ID          int                    `json:"id"`
	SKU         string                 `json:"sku"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	CategoryID  string                 `json:"categoryId"`
	IsActive    bool                   `json:"isActive"`
	ImageURL    string                 `json:"imageUrl"`
	Price       float64                `json:"price"`
	StockQty    int                    `json:"stockQty"`
	Attributes  map[string]interface{} `json:"attributes"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}
