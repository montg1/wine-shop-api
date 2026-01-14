package domain

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID uint        `json:"user_id"`
	Total  float64     `json:"total"`
	Status string      `json:"status"` // pending, paid, shipped, cancelled
	Items  []OrderItem `json:"items"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // Price at time of purchase
}
