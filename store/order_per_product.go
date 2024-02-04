package store

import (
	"github.com/NikhilSharmaWe/market/models"
	"gorm.io/gorm"
)

type OrderPerProductStore interface {
	CreateTable() error
	Create(fr models.OrderPerProductDBModel) error
	Update(updateMap, whereMap map[string]interface{}) error
	Delete(whereMap map[string]interface{}) error
	GetMany(whereMap map[string]interface{}) ([]models.OrderPerProductDBModel, error)
	IsExists(whereMap map[string]interface{}) (bool, error)
	DB() *gorm.DB
}

type orderPerProductStore struct {
	db *gorm.DB
}

func NewOrderPerProductStore(db *gorm.DB) OrderPerProductStore {
	return &orderPerProductStore{
		db: db,
	}
}

func (os *orderPerProductStore) table() string {
	return "order_per_product"
}

func (os *orderPerProductStore) DB() *gorm.DB {
	return os.db
}

func (os *orderPerProductStore) CreateTable() error {
	return os.db.Table(os.table()).AutoMigrate(models.OrderPerProductDBModel{})

}

func (os *orderPerProductStore) Create(fr models.OrderPerProductDBModel) error {
	return os.db.Table(os.table()).Create(fr).Error
}

func (os *orderPerProductStore) Update(updateMap, whereMap map[string]interface{}) error {
	return os.db.Table(os.table()).Where(whereMap).Updates(updateMap).Error
}

func (os *orderPerProductStore) Delete(whereMap map[string]interface{}) error {
	return os.db.Table(os.table()).Where(whereMap).Delete(nil).Error
}

func (os *orderPerProductStore) GetMany(whereMap map[string]interface{}) ([]models.OrderPerProductDBModel, error) {
	resp := []models.OrderPerProductDBModel{}
	if err := os.db.Table(os.table()).Where(whereMap).Find(&resp).Error; err != nil {
		return resp, err
	}

	if len(resp) == 0 {
		return resp, models.ErrMatchingRecordNotFound
	}

	return resp, nil
}

func (os *orderPerProductStore) IsExists(whereMap map[string]interface{}) (bool, error) {
	var count int64
	err := os.db.Table(os.table()).Where(whereMap).Count(&count).Error
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
