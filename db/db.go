package db

import (
	"fmt"

	"github.com/anhdt-vnpay/f5_fulltext_search/lib/config"
	m "github.com/anhdt-vnpay/f5_fulltext_search/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func InitMysql() (*gorm.DB, error) {
	fmt.Println("init mysql database ...")
	config := config.GetConfig()

	// Declare variable configuration mysql
	host := config.GetString("mysql.host")
	port := config.GetString("mysql.port")
	dbName := config.GetString("mysql.database")
	user := config.GetString("mysql.username")
	pass := config.GetString("mysql.password")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&m.User{},
	)
	return db, err
}

func InitPostgres() (*gorm.DB, error) {
	fmt.Println("init postgres database ...")
	config := config.GetConfig()

	// Declare variable configuration mysql
	host := config.GetString("mysql.host")
	port := config.GetString("mysql.port")
	dbName := config.GetString("mysql.database")
	user := config.GetString("mysql.username")
	pass := config.GetString("mysql.password")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh", host, user, pass, dbName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&m.User{},
	)
	return db, err
}
