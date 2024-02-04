package store

import (
	"github.com/NikhilSharmaWe/market/models"
	"gorm.io/gorm"
)

type ProductStore interface {
	CreateTable() error
	Create(fr models.ProductDBModel) error
	Update(updateMap, whereMap map[string]interface{}) error
	Delete(whereMap map[string]interface{}) error
	GetOne(whereMap map[string]interface{}) (*models.ProductDBModel, error)
	GetMany(whereMap map[string]interface{}) ([]models.ProductDBModel, error)
	IsExists(whereMap map[string]interface{}) (bool, error)
	DB() *gorm.DB
}

type productStore struct {
	db *gorm.DB
}

func NewProductsStore(db *gorm.DB) ProductStore {
	return &productStore{
		db: db,
	}
}

func (ps *productStore) table() string {
	return "products"
}

func (ps *productStore) DB() *gorm.DB {
	return ps.db
}

func (ps *productStore) CreateTable() error {
	return ps.db.Table(ps.table()).AutoMigrate(models.ProductDBModel{})
}

func (ps *productStore) Create(fr models.ProductDBModel) error {
	return ps.db.Table(ps.table()).Create(fr).Error
}

func (ps *productStore) Update(updateMap, whereMap map[string]interface{}) error {
	return ps.db.Table(ps.table()).Where(whereMap).Updates(updateMap).Error
}

func (ps *productStore) Delete(whereMap map[string]interface{}) error {
	return ps.db.Table(ps.table()).Where(whereMap).Delete(nil).Error
}

func (ps *productStore) GetOne(whereMap map[string]interface{}) (*models.ProductDBModel, error) {
	var product models.ProductDBModel
	if err := ps.db.Table(ps.table()).Where(whereMap).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (os *productStore) GetMany(whereMap map[string]interface{}) ([]models.ProductDBModel, error) {
	resp := []models.ProductDBModel{}
	if err := os.db.Table(os.table()).Where(whereMap).Find(&resp).Error; err != nil {
		return resp, err
	}

	if len(resp) == 0 {
		return resp, models.ErrMatchingRecordNotFound
	}

	return resp, nil
}

func (ps *productStore) IsExists(whereMap map[string]interface{}) (bool, error) {
	var count int64
	err := ps.db.Table(ps.table()).Where(whereMap).Count(&count).Error
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
