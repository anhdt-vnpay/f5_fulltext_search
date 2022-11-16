package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint64     `json:"id"`
	Name     string     `json:"name" gorm:"unique"`
	Type     string     `json:"type"`
	CreateAt *time.Time `json:"create_at" gorm:"type:timestamp;autoCreateTime"`
	UpdateAt *time.Time `json:"update_at" gorm:"type:timestamp"`
	Delete   bool       `json:"delete"`
}

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (u *UserHandler) List() ([]User, error) {
	var users []User
	res := u.db.Model(User{}).Find(&users)
	if res.Error != nil {
		return users, res.Error
	}
	return users, nil
}

func (u *UserHandler) Detail(user User) (User, error) {
	var ru User
	res := u.db.Model(User{}).Where("id = ?", user.ID).Find(&ru)
	if res.Error != nil {
		return User{}, res.Error
	}
	return ru, nil
}

func (u *UserHandler) Create(user User) (User, error) {
	var ru User
	if u.db.Model(User{}).Where("name", user.Name).First(&ru).RowsAffected > 0 {
		return User{}, fmt.Errorf("User existed %s", user.Name)
	}

	if res := u.db.Create(&user); res.Error != nil {
		return User{}, res.Error
	}
	return user, nil
}

func (u *UserHandler) Update(user User) (User, error) {
	res := u.db.Model(User{}).Where("id = ?", user.ID).Updates(&user)
	if res.Error != nil {
		return User{}, res.Error
	}
	return user, nil
}

func (u *UserHandler) Delete(user User) (User, error) {
	res := u.db.Model(User{}).Where("id = ?", user.ID).Delete(&user)
	if res.Error != nil {
		return user, res.Error
	}
	return user, nil
}
