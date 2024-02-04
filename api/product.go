package api

import (
	"github.com/NikhilSharmaWe/market/models"
	"github.com/NikhilSharmaWe/market/store"
)

type ProductService interface {
	Create(*models.CreateProductRequest) error
	Get(*models.GetProductsRequest) (map[string]models.ProductResponse, error)
	Delete(*models.DeleteProductRequest) error
	Update(*models.UpdateProductRequest) error
}

type productService struct {
	store.ProductStore
	store.InventoryStore
}

func NewProductService(productStore store.ProductStore, inventoryStore store.InventoryStore) ProductService {
	return &productService{
		ProductStore:   productStore,
		InventoryStore: inventoryStore,
	}
}

func (ps *productService) Create(req *models.CreateProductRequest) error {
	exists, err := ps.ProductStore.IsExists(map[string]interface{}{"product_name": req.ProductName})
	if err != nil {
		return err
	}

	if exists {
		return models.ErrProductAlreadyExists
	}

	if err := ps.ProductStore.Create(models.ProductDBModel{
		ProductName: req.ProductName,
		Category:    req.Category,
		Price:       req.Price,
	}); err != nil {
		return err
	}

	if err := ps.InventoryStore.Create(models.InventoryDBModel{
		ProductName:     req.ProductName,
		QuantityInStock: req.InitialQuantity,
	}); err != nil {
		return err
	}

	return nil
}

func (ps *productService) Get(req *models.GetProductsRequest) (map[string]models.ProductResponse, error) {

	var wheremap map[string]interface{}
	if req == nil {
		wheremap = nil
	} else {
		wheremap = map[string]interface{}{"category": req.Category}
	}

	products, err := ps.ProductStore.GetMany(wheremap)
	if err != nil {
		return nil, err
	}

	resp := map[string]models.ProductResponse{}
	for _, product := range products {
		resp[product.ProductName] = models.ProductResponse{
			Category: product.Category,
			Price:    product.Price,
		}
	}

	return resp, nil
}

func (ps *productService) Delete(req *models.DeleteProductRequest) error {
	exists, err := ps.ProductStore.IsExists(map[string]interface{}{"product_name": req.ProductName})
	if err != nil {
		return err
	}

	if !exists {
		return models.ErrProductNotFound
	}

	return ps.ProductStore.Delete(map[string]interface{}{"product_name": req.ProductName})
}

func (ps *productService) Update(req *models.UpdateProductRequest) error {
	exists, err := ps.ProductStore.IsExists(map[string]interface{}{"product_name": req.ProductName})
	if err != nil {
		return err
	}

	if !exists {
		return models.ErrProductNotFound
	}

	updateMap := map[string]interface{}{}
	if req.Category == "" && req.Price == 0 {
		return nil
	}

	if req.Category != "" {
		updateMap["category"] = req.Category
	}
	if req.Price != 0 {
		updateMap["price"] = req.Price
	}

	return ps.ProductStore.Update(updateMap, map[string]interface{}{"product_name": req.ProductName})
}
