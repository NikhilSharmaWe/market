package store

import (
	"github.com/NikhilSharmaWe/market/models"
	"gorm.io/gorm"
)

type AdminsStore interface {
	CreateTable() error
	Create(fr models.AdminDBModel) error
	Update(updateMap, whereMap map[string]interface{}) error
	Delete(whereMap map[string]interface{}) error
	IsExists(whereMap map[string]interface{}) (bool, error)
	DB() *gorm.DB
}

type adminStore struct {
	db *gorm.DB
}

func NewAdminsStore(db *gorm.DB) AdminsStore {
	return &adminStore{
		db: db,
	}
}

func (as *adminStore) table() string {
	return "admins"
}

func (as *adminStore) DB() *gorm.DB {
	return as.db
}

func (as *adminStore) CreateTable() error {
	return as.db.Table(as.table()).AutoMigrate(models.AdminDBModel{})
}

func (as *adminStore) Create(fr models.AdminDBModel) error {
	return as.db.Table(as.table()).Create(fr).Error
}

func (as *adminStore) Update(updateMap, whereMap map[string]interface{}) error {
	return as.db.Table(as.table()).Where(whereMap).Updates(updateMap).Error
}

func (as *adminStore) Delete(whereMap map[string]interface{}) error {
	return as.db.Table(as.table()).Where(whereMap).Delete(nil).Error
}

func (as *adminStore) IsExists(whereMap map[string]interface{}) (bool, error) {
	var count int64
	err := as.db.Table(as.table()).Where(whereMap).Count(&count).Error
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
