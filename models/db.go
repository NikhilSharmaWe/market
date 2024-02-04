package models

import "time"

type ProductDBModel struct {
	ProductName string `gorm:"column:product_name;primaryKey"`
	Category    string `gorm:"column:category"`
	Price       int    `gorm:"column:price"`
}

type UserDBModel struct {
	Username string `gorm:"column:username;primaryKey"`
	Email    string `gorm:"column:email"`
	Password []byte `gorm:"column:password_hash"`
}

type AdminDBModel struct {
	Username string `gorm:"column:username;primaryKey"`
}

type OrderDBModel struct {
	OrderID   string    `gorm:"column:order_id;primaryKey"`
	Username  string    `gorm:"column:username"`
	TotalCost int       `gorm:"column:total_cost"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type OrderPerProductDBModel struct {
	OrderID     string `gorm:"column:order_id" json:"created_at"`
	ProductName string `gorm:"column:product_name" json:"product_name"`
	Quantity    int    `gorm:"column:quantity" json:"quantity"`
}

type InventoryDBModel struct {
	ProductName     string `gorm:"column:product_name;primaryKey"`
	QuantityInStock int    `gorm:"column:quantity_in_stock"`
}
