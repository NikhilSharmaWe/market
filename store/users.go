package store

import (
	"github.com/NikhilSharmaWe/market/models"
	"gorm.io/gorm"
)

type UsersStore interface {
	CreateTable() error
	Create(fr models.UserDBModel) error
	Update(updateMap, whereMap map[string]interface{}) error
	Delete(whereMap map[string]interface{}) error
	GetOne(whereMap map[string]interface{}) (*models.UserDBModel, error)
	IsExists(whereMap map[string]interface{}) (bool, error)
	DB() *gorm.DB
}

type userStore struct {
	db *gorm.DB
}

func NewUsersStore(db *gorm.DB) UsersStore {
	return &userStore{
		db: db,
	}
}

func (us *userStore) table() string {
	return "users"
}

func (us *userStore) DB() *gorm.DB {
	return us.db
}

func (us *userStore) CreateTable() error {
	return us.db.Table(us.table()).AutoMigrate(models.UserDBModel{})
}

func (us *userStore) Create(fr models.UserDBModel) error {
	return us.db.Table(us.table()).Create(fr).Error
}

func (us *userStore) Update(updateMap, whereMap map[string]interface{}) error {
	return us.db.Table(us.table()).Where(whereMap).Updates(updateMap).Error
}

func (us *userStore) Delete(whereMap map[string]interface{}) error {
	return us.db.Table(us.table()).Where(whereMap).Delete(nil).Error
}

func (us *userStore) GetOne(whereMap map[string]interface{}) (*models.UserDBModel, error) {
	var user models.UserDBModel
	if err := us.db.Table(us.table()).Where(whereMap).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *userStore) IsExists(whereMap map[string]interface{}) (bool, error) {
	var count int64
	err := us.db.Table(us.table()).Where(whereMap).Count(&count).Error
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
