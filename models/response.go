package models

import "time"

type OrdersResponse struct {
	OrderID    string    `json:"order_id"`
	TotalCost  int       `json:"total_cost"`
	Products   []string  `json:"products"`
	Quantities []int     `json:"quantities"`
	CreatedAt  time.Time `json:"created_at"`
}

type ProductResponse struct {
	Category string `json:"category"`
	Price    int    `json:"price"`
}
