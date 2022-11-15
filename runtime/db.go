package runtime

import (
	"fmt"

	helper "github.com/anhdt-vnpay/f5_fulltext_search/helper"
	m "github.com/anhdt-vnpay/f5_fulltext_search/model"
	"gorm.io/gorm"
)

type DbStorage interface {
	Insert(tableName string, value interface{}) error
	Update(tableName string, value interface{}) error
	Delete(tableName string, value interface{}) error
	Get(tableName string, condition string) error
}

type gormDbStorage struct {
	db *gorm.DB
}

func NewGormDbStorage(db *gorm.DB) *gormDbStorage {
	return &gormDbStorage{
		db: db,
	}
}

func (r *gormDbStorage) Insert(tableName string, value interface{}) error {
	model, err := helper.ModelByName(tableName)
	if err != nil {
		return err
	}

	user, ok := value.(m.User)
	if !ok {
		fmt.Println("convert failed")
		return fmt.Errorf("convert failed")
	}

	var exist m.User
	if res := r.db.Model(model).Where("name = ?", user.Name).First(&exist); res.RowsAffected > 0 {
		return fmt.Errorf("user existed")
	}

	err = r.db.Model(model).Create(&user).Error
	return err
}

func (r *gormDbStorage) Update(tableName string, value interface{}) error {
	model, err := helper.ModelByName(tableName)
	if err != nil {
		return err
	}

	user, ok := value.(m.User)
	if !ok {
		fmt.Println("convert failed")
		return fmt.Errorf("convert failed")
	}

	var exist m.User
	if res := r.db.Model(model).Where("id = ?", user.ID).First(&exist); res.RowsAffected <= 0 {
		fmt.Println("Create user instead")
		err := r.db.Model(model).Create(&user).Error
		return err
	}

	err = r.db.Model(model).Where("id = ?", user.ID).Updates(&user).Error
	return err
}

func (r *gormDbStorage) Delete(tableName string, value interface{}) error {
	model, err := helper.ModelByName(tableName)
	if err != nil {
		return err
	}

	user, ok := value.(m.User)
	if !ok {
		fmt.Println("convert failed")
		return fmt.Errorf("convert failed")
	}

	err = r.db.Model(model).Delete(&user).Error
	return err
}

func (r *gormDbStorage) Get(tableName string, condition string) error {
	fmt.Println("Get gorm >>>")
	return nil
}
