package runtime

import (
	"fmt"
)

type DbFullTextSearch interface {
	Insert(tableName string, value interface{}) error
	Update(tableName string, value interface{}) error
	Delete(tableName string, value interface{}) error
	Get(tableName string, condition string) error
	SearchLite(query string, result interface{}) error

	SetStorage(dbStorage DbStorage)
	SetMsgProcessor(msgProcessor MessageProcessor)
	SetSearchProcessor(search SearchProcessor)
}

type DbOption func(db DbFullTextSearch)

type dbFullTextSearch struct {
	db           DbStorage
	msgProcessor MessageProcessor
	search       SearchProcessor
}

func NewDbFullTextSearch(opts ...DbOption) DbFullTextSearch {
	db := &dbFullTextSearch{}
	for _, opt := range opts {
		opt(db)
	}
	return db
}

func (db *dbFullTextSearch) Insert(tableName string, value interface{}) error {
	if err := db.db.Insert(tableName, value); err != nil {
		return err
	}
	return nil
}

func (db *dbFullTextSearch) Update(tableName string, value interface{}) error {
	if err := db.db.Update(tableName, value); err != nil {
		return err
	}
	return nil
}

func (db *dbFullTextSearch) Delete(tableName string, value interface{}) error {
	if err := db.db.Delete(tableName, value); err != nil {
		return err
	}
	return nil
}
func (db *dbFullTextSearch) Get(tableName string, condition string) error {
	fmt.Println("Get >>")
	err := db.db.Get(tableName, condition)
	return err
}

func (db *dbFullTextSearch) SearchLite(query string, result interface{}) error {
	fmt.Println("Search lite >>")
	_, err := db.search.SearchLite(query)
	if err != nil {
		fmt.Println("Search lite error: ", err.Error())
		return err
	}
	return nil
}

func (db *dbFullTextSearch) SetStorage(dbStorage DbStorage) {
	fmt.Println("Set storage >>")
	db.db = dbStorage
}

func (db *dbFullTextSearch) SetMsgProcessor(msgProcessor MessageProcessor) {
	fmt.Println("Set msg processor >>")
	db.msgProcessor = msgProcessor
}

func (db *dbFullTextSearch) SetSearchProcessor(search SearchProcessor) {
	fmt.Println("Set search processor >>")
	db.search = search
}

/*************************************************************************************************
With Option
*************************************************************************************************/
func WithStorage(dbStorage DbStorage) DbOption {
	return func(db DbFullTextSearch) {
		db.SetStorage(dbStorage)
	}
}

func WithMsgProcessor(msgProcessor MessageProcessor) DbOption {
	return func(db DbFullTextSearch) {
		db.SetMsgProcessor(msgProcessor)
	}
}

func WithSearchProcessor(search SearchProcessor) DbOption {
	return func(db DbFullTextSearch) {
		db.SetSearchProcessor(search)
	}
}
