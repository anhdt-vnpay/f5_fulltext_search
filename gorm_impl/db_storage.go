package db_storage

import (
	"gorm.io/gorm"
)

type gormDbStorage struct {
	db *gorm.DB
}

func NewGormDbStorage(db *gorm.DB) *gormDbStorage {
	return &gormDbStorage{
		db: db,
	}

}

func (r *gormDbStorage) Insert(tableName string, value interface{}) error {
	err := r.db.Table(tableName).Create(value).Error
	return err
}

func (r *gormDbStorage) Update(tableName string, value interface{}) error {
	err := r.db.Table(tableName).Updates(value).Error
	return err
}

func (r *gormDbStorage) Delete(tableName string, value interface{}) error {
	err := r.db.Table(tableName).Delete(value).Error
	return err
}

func (r *gormDbStorage) Get(tableName string, condition string, result interface{}) error {
	err := r.db.Table(tableName).Where(condition).Find(result).Error
	return err
}
