package store

import (
	"github.com/NikhilSharmaWe/market/models"
	"gorm.io/gorm"
)

type OrderStore interface {
	CreateTable() error
	Create(fr models.OrderDBModel) error
	Update(updateMap, whereMap map[string]interface{}) error
	Delete(whereMap map[string]interface{}) error
	GetOne(whereMap map[string]interface{}) (*models.OrderDBModel, error)
	GetMany(whereMap map[string]interface{}, orderBy string) ([]models.OrderDBModel, error)
	DB() *gorm.DB
}

type orderStore struct {
	db *gorm.DB
}

func NewOrdersStore(db *gorm.DB) OrderStore {
	return &orderStore{
		db: db,
	}
}

func (os *orderStore) table() string {
	return "orders"
}

func (os *orderStore) DB() *gorm.DB {
	return os.db
}

func (os *orderStore) CreateTable() error {
	return os.db.Table(os.table()).AutoMigrate(models.OrderDBModel{})
}

func (os *orderStore) Create(fr models.OrderDBModel) error {
	return os.db.Table(os.table()).Create(fr).Error
}

func (os *orderStore) Update(updateMap, whereMap map[string]interface{}) error {
	return os.db.Table(os.table()).Where(whereMap).Updates(updateMap).Error
}

func (os *orderStore) Delete(whereMap map[string]interface{}) error {
	return os.db.Table(os.table()).Where(whereMap).Delete(nil).Error
}

func (os *orderStore) GetOne(whereMap map[string]interface{}) (*models.OrderDBModel, error) {
	var order models.OrderDBModel
	if err := os.db.Table(os.table()).Where(whereMap).First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (os *orderStore) GetMany(whereMap map[string]interface{}, orderBy string) ([]models.OrderDBModel, error) {
	resp := []models.OrderDBModel{}
	if err := os.db.Table(os.table()).Order(orderBy).Where(whereMap).Find(&resp).Error; err != nil {
		return resp, err
	}

	if len(resp) == 0 {
		return resp, models.ErrMatchingRecordNotFound
	}

	return resp, nil
}
