package api

import (
	"fmt"
	"time"

	"github.com/NikhilSharmaWe/market/models"
	"github.com/NikhilSharmaWe/market/store"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	Create(*models.CreateOrderRequest) (*models.OrderDBModel, error)
	GetOrdersByUsername(username string) (map[string]models.OrdersResponse, error)
	GetOrdersForAdmin(*models.GetOrdersForAdminRequest) (map[string]models.OrdersResponse, error)
}

type orderService struct {
	store.OrderStore
	store.OrderPerProductStore
	store.ProductStore
	store.InventoryStore
}

func NewOrderService(orderStore store.OrderStore, orderPerProductStore store.OrderPerProductStore, productStore store.ProductStore, inventoryStore store.InventoryStore) OrderService {
	return &orderService{
		OrderStore:           orderStore,
		OrderPerProductStore: orderPerProductStore,
		ProductStore:         productStore,
		InventoryStore:       inventoryStore,
	}
}

func (os *orderService) Create(req *models.CreateOrderRequest) (*models.OrderDBModel, error) {
	var order models.OrderDBModel
	db := os.InventoryStore.DB()

	if err := db.Transaction(func(tx *gorm.DB) error {
		var totalCost int
		productStore := store.NewProductsStore(tx)
		inventoryStore := store.NewInventoryStore(tx)
		orderStore := store.NewOrdersStore(tx)
		orderPerProductStore := store.NewOrderPerProductStore(tx)
		orderID := uuid.NewV4().String()

		order = models.OrderDBModel{
			OrderID:   orderID,
			Username:  req.Username,
			TotalCost: 0,
			CreatedAt: time.Now(),
		}

		if err := orderStore.Create(order); err != nil {
			return err
		}

		for productName, quantity := range req.ProductsAndQuantity {
			exists, err := productStore.IsExists(map[string]interface{}{"product_name": productName})
			if err != nil {
				return err
			}

			if !exists {
				return models.ErrProductNotFound
			}

			product, err := productStore.GetOne(map[string]interface{}{"product_name": productName})
			if err != nil {
				return err
			}

			inventory, err := inventoryStore.GetOne(map[string]interface{}{"product_name": productName})
			if err != nil {
				return err
			}

			if inventory.QuantityInStock < quantity {
				return models.ErrInvalidQuantity
			}

			totalCost += product.Price * quantity
			order.TotalCost = totalCost

			uq := inventory.QuantityInStock - quantity

			if err := orderPerProductStore.Create(models.OrderPerProductDBModel{
				OrderID:     orderID,
				ProductName: productName,
				Quantity:    quantity,
			}); err != nil {
				return err
			}

			inventory.QuantityInStock = uq
			if err := inventoryStore.Update(map[string]interface{}{
				"quantity_in_stock": uq,
			}, map[string]interface{}{"product_name": productName}); err != nil {
				return err
			}
		}

		return orderStore.Update(map[string]interface{}{
			"total_cost": totalCost,
		}, map[string]interface{}{
			"order_id": orderID,
		})
	}); err != nil {
		return nil, err
	}

	return &order, nil
}

func (os *orderService) GetOrdersByUsername(username string) (map[string]models.OrdersResponse, error) {
	resp := map[string]models.OrdersResponse{}
	orders, err := os.OrderStore.GetMany(map[string]interface{}{"username": username}, "created_at")
	if err != nil {
		return nil, err
	}

	for _, order := range orders {

		ordersPerProduct, err := os.OrderPerProductStore.GetMany(map[string]interface{}{"order_id": order.OrderID})
		if err != nil {
			return nil, err
		}

		productsName := []string{}
		quantities := []int{}

		for _, opp := range ordersPerProduct {
			productsName = append(productsName, opp.ProductName)
			quantities = append(quantities, opp.Quantity)
		}

		resp[order.OrderID] = models.OrdersResponse{
			OrderID:    order.OrderID,
			TotalCost:  order.TotalCost,
			Products:   productsName,
			Quantities: quantities,
			CreatedAt:  order.CreatedAt,
		}
	}

	return resp, nil
}

func (os *orderService) GetOrdersForAdmin(req *models.GetOrdersForAdminRequest) (map[string]models.OrdersResponse, error) {
	resp := map[string]models.OrdersResponse{}

	query := os.OrderStore.DB().Table("orders")

	if req.Username != "" {
		query = query.Where("username = ?", req.Username)
	}

	if req.ProductName != "" {
		query = query.Joins("JOIN order_per_product ON orders.order_id = order_per_product.order_id").
			Where("order_per_product.product_name = ?", req.ProductName)
	}

	var orders []models.OrderDBModel
	if err := query.Order(req.SortBy).Find(&orders).Error; err != nil {
		return nil, err
	}

	for i, order := range orders {
		ordersPerProduct, err := os.OrderPerProductStore.GetMany(map[string]interface{}{"order_id": order.OrderID})
		if err != nil {
			return nil, err
		}

		productsName := []string{}
		quantities := []int{}

		for _, opp := range ordersPerProduct {
			productsName = append(productsName, opp.ProductName)
			quantities = append(quantities, opp.Quantity)
		}

		resp[fmt.Sprint(i)+"__"+order.OrderID+"__"+order.Username] = models.OrdersResponse{
			OrderID:    order.OrderID,
			TotalCost:  order.TotalCost,
			Products:   productsName,
			Quantities: quantities,
			CreatedAt:  order.CreatedAt,
		}
	}

	return resp, nil
}
