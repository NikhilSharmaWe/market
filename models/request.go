package models

type CreateProductRequest struct {
	ProductName     string `json:"product_name"`
	Category        string `json:"category"`
	Price           int    `json:"price"`
	InitialQuantity int    `json:"initial_quantity"`
}

type DeleteProductRequest struct {
	ProductName string `json:"product_name"`
}

type UpdateProductRequest struct {
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
	Price       int    `json:"price"`
}

type UpdateInventoryRequest struct {
	ProductName string `json:"product_name"`
	Operation   string `json:"operation"`
	Quantity    int    `json:"quantity"`
}

type CreateOrderRequest struct {
	Username            string         `json:"username"`
	ProductsAndQuantity map[string]int `json:"product_and_quantity"`
}

type GetOrdersForAdminRequest struct {
	Username    string `json:"username"`
	ProductName string `json:"product_name"`
	SortBy      string `json:"sort_by"`
}

type GetProductsRequest struct {
	Category string `json:"category"`
}
