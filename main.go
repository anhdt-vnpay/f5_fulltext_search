package main

import (
	"fmt"
	"log"

	gorm_impl "github.com/anhdt-vnpay/f5_fulltext_search/gorm_impl"
	db_connector "github.com/anhdt-vnpay/f5_fulltext_search/gorm_impl/db_connector"
	"github.com/anhdt-vnpay/f5_fulltext_search/lib/config"
	m "github.com/anhdt-vnpay/f5_fulltext_search/model"
	"github.com/anhdt-vnpay/f5_fulltext_search/runtime"
)

func main() {
	fmt.Println("====================== Init DB ==========================")
	config.Init()
	fmt.Println("init mysql database ...")
	config := config.GetConfig()

	// Declare variable configuration mysql
	host := config.GetString("mysql.host")
	port := config.GetString("mysql.port")
	dbName := config.GetString("mysql.database")
	user := config.GetString("mysql.username")
	pass := config.GetString("mysql.password")

	configObj := db_connector.ConnectorConfig{
		Mode:     db_connector.Mysql,
		Host:     host,
		Port:     port,
		Database: dbName,
		Username: user,
		Password: pass,
	}

	conn, err := db_connector.NewDatabaseConnector(&configObj)
	if err != nil {
		log.Fatal("DB connector error: ", err.Error())
	}

	db := conn.GetDb()

	fmt.Println("====================== Init objects ==========================")

	dbStorage := gorm_impl.NewGormDbStorage(db)
	var opts []runtime.DbOption
	opt1 := runtime.WithStorage(dbStorage)
	opt2 := runtime.WithMsgProcessor(nil)
	opts = append(opts, opt1, opt2)
	dbf := runtime.NewDbFullTextSearch(opts...)

	fmt.Println("====================== DEMO ==========================")

	// Test get
	fmt.Println("Get >>>>>>>>")
	rs := []m.User{}
	name := "Rename A"
	err = dbf.Get("users", fmt.Sprintf("name=%q", name), &rs)
	if err != nil {
		fmt.Println("Get error: ", err.Error())
	}
	fmt.Println("RS: ")
	fmt.Println(rs)
}
