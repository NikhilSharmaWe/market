package api

import (
	"github.com/NikhilSharmaWe/market/models"
	"github.com/NikhilSharmaWe/market/store"
)

type InventoryService interface {
	Update(*models.UpdateInventoryRequest) error
}

type inventoryService struct {
	store.ProductStore
	store.InventoryStore
}

func NewInventoryService(productStore store.ProductStore, inventoryStore store.InventoryStore) InventoryService {
	return &inventoryService{
		ProductStore:   productStore,
		InventoryStore: inventoryStore,
	}
}

func (is *inventoryService) Update(req *models.UpdateInventoryRequest) error {
	if req.Operation != "add" && req.Operation != "remove" {
		return models.ErrInvalidOperaton
	}

	exists, err := is.ProductStore.IsExists(map[string]interface{}{"product_name": req.ProductName})
	if err != nil {
		return err
	}

	if !exists {
		return models.ErrProductNotFound
	}

	inventory, err := is.InventoryStore.GetOne(map[string]interface{}{"product_name": req.ProductName})
	if err != nil {
		return err
	}

	var uq int

	if req.Operation == "add" {
		uq = inventory.QuantityInStock + req.Quantity
	} else {
		uq = inventory.QuantityInStock - req.Quantity
		if uq < 0 {
			return models.ErrInvalidQuantity
		}
	}

	inventory.QuantityInStock = uq
	if err := is.InventoryStore.Update(map[string]interface{}{
		"quantity_in_stock": uq,
	}, map[string]interface{}{"product_name": req.ProductName}); err != nil {
		return err
	}

	return nil
}
