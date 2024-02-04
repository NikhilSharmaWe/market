package store

import (
	"github.com/NikhilSharmaWe/market/models"
	"gorm.io/gorm"
)

type InventoryStore interface {
	CreateTable() error
	Create(fr models.InventoryDBModel) error
	Update(updateMap, whereMap map[string]interface{}) error
	Delete(whereMap map[string]interface{}) error
	GetOne(whereMap map[string]interface{}) (*models.InventoryDBModel, error)
	DB() *gorm.DB
}

type inventoryStore struct {
	db *gorm.DB
}

func NewInventoryStore(db *gorm.DB) InventoryStore {
	return &inventoryStore{
		db: db,
	}
}

func (is *inventoryStore) table() string {
	return "inventory"
}

func (is *inventoryStore) DB() *gorm.DB {
	return is.db
}

func (is *inventoryStore) CreateTable() error {
	return is.db.Table(is.table()).AutoMigrate(models.InventoryDBModel{})

}

func (is *inventoryStore) Create(fr models.InventoryDBModel) error {
	return is.db.Table(is.table()).Create(fr).Error
}

func (is *inventoryStore) Update(updateMap, whereMap map[string]interface{}) error {
	return is.db.Table(is.table()).Where(whereMap).Updates(updateMap).Error
}

func (is *inventoryStore) Delete(whereMap map[string]interface{}) error {
	return is.db.Table(is.table()).Where(whereMap).Delete(nil).Error
}

func (is *inventoryStore) GetOne(whereMap map[string]interface{}) (*models.InventoryDBModel, error) {
	var inventory models.InventoryDBModel
	if err := is.db.Table(is.table()).Where(whereMap).First(&inventory).Error; err != nil {
		return nil, err
	}

	return &inventory, nil
}
