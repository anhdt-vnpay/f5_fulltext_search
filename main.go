package main

import (
	"fmt"
	"log"

	"github.com/anhdt-vnpay/f5_fulltext_search/db"
	"github.com/anhdt-vnpay/f5_fulltext_search/lib/config"
	m "github.com/anhdt-vnpay/f5_fulltext_search/model"
	runtime "github.com/anhdt-vnpay/f5_fulltext_search/runtime"
)

func main() {
	fmt.Println("====================== INIT DB ==========================")

	config.Init()
	db, err := db.InitMysql()
	if err != nil {
		log.Fatal("Init db failed: ", err.Error())
	}

	fmt.Println("====================== Init objects ==========================")

	dbStorage := runtime.NewGormDbStorage(db)
	var opts []runtime.DbOption
	opt1 := runtime.WithStorage(dbStorage)
	opt2 := runtime.WithMsgProcessor(nil)
	opts = append(opts, opt1, opt2)
	dbf := runtime.NewDbFullTextSearch(opts...)

	fmt.Println("====================== DEMO ==========================")

	// Test insert
	fmt.Println("Insert >>>>>>>>>>>>")
	a := m.User{
		ID:   1,
		Name: "A",
		Type: "Person",
	}

	b := m.User{
		ID:   2,
		Name: "B",
		Type: "Person",
	}

	c := m.User{
		ID:   3,
		Name: "C",
		Type: "Person",
	}

	if err := dbf.Insert("user", a); err != nil {
		fmt.Println("[Insert] error occured: ", err.Error())
	}
	if err := dbf.Insert("user", b); err != nil {
		fmt.Println("[Insert] error occured: ", err.Error())
	}
	if err := dbf.Insert("user", c); err != nil {
		fmt.Println("[Insert] error occured: ", err.Error())
	}
	fmt.Println("[Insert] SUCCESS")
	fmt.Println("")

	// Test update
	fmt.Println("Update >>>>>>>>>>>>")
	editA := m.User{
		ID:   1,
		Name: "Rename A",
		Type: "Person",
	}
	if err := dbf.Update("user", editA); err != nil {
		fmt.Println("[Update] error occured: ", err.Error())
	}
	fmt.Println("[Update] SUCCESS")
	fmt.Println("")

	// Test delete
	fmt.Println("Delete >>>>>>>>>>>>")
	deleteC := m.User{
		ID:   3,
		Name: "C",
		Type: "Person",
	}
	if err := dbf.Delete("user", deleteC); err != nil {
		fmt.Println("[Delete] error occured: ", err.Error())
	}
	fmt.Println("[Delete] SUCCESS")
	fmt.Println("")
}
