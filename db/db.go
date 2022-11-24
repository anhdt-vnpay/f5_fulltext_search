package db

import (
	"fmt"

	"github.com/anhdt-vnpay/f5_fulltext_search/lib/config"
	"github.com/anhdt-vnpay/f5_fulltext_search/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func Init() (*gorm.DB, error) {
	fmt.Println("Init database ...")
	config := config.GetConfig()

	// Declare variables configuration mysql
	host := config.GetString("mysql.host")
	port := config.GetString("mysql.port")
	dbName := config.GetString("mysql.database")
	user := config.GetString("mysql.username")
	pass := config.GetString("mysql.password")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbName)
	fmt.Printf("db query: %s \n", dsn)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&model.User{},
	)

	return db, err
}
