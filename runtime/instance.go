package runtime

import (
	"fmt"
)

var (
	TYPE_INSERT = "insert"
	TYPE_UPDATE = "update"
	TYPE_DELETE = "delete"
)

type DbFullTextSearch interface {
	Insert(tableName string, value interface{}) error
	Update(tableName string, value interface{}) error
	Delete(tableName string, value interface{}) error
	Get(tableName string, condition string, result interface{}) error
	SearchLite(query string, result interface{}) error
}

type DbOption func(db *dbFullTextSearch)

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
	go db.msgLoop()
	return db
}

func (db *dbFullTextSearch) Insert(tableName string, value interface{}) error {

	if err := db.db.Insert(tableName, value); err != nil {
		return err
	}
	go db.saveMsg(TYPE_INSERT, value)
	return nil
}

func (db *dbFullTextSearch) Update(tableName string, value interface{}) error {
	if err := db.db.Update(tableName, value); err != nil {
		return err
	}
	go db.saveMsg(TYPE_UPDATE, value)
	return nil
}

func (db *dbFullTextSearch) Delete(tableName string, value interface{}) error {
	if err := db.db.Delete(tableName, value); err != nil {
		return err
	}
	go db.saveMsg(TYPE_DELETE, value)
	return nil
}
func (db *dbFullTextSearch) Get(tableName string, condition string, result interface{}) error {
	fmt.Println("Get >>")
	err := db.db.Get(tableName, condition, result)
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

/*************************************************************************************************
With Option
*************************************************************************************************/
func WithStorage(dbStorage DbStorage) DbOption {
	return func(db *dbFullTextSearch) {
		db.db = dbStorage
	}
}

func WithMsgProcessor(msgProcessor MessageProcessor) DbOption {
	return func(db *dbFullTextSearch) {
		db.msgProcessor = msgProcessor
	}
}

func WithSearchProcessor(search SearchProcessor) DbOption {
	return func(db *dbFullTextSearch) {
		db.search = search
	}
}

func (db *dbFullTextSearch) saveMsg(tipe string, data interface{}) {
	message, err := createMessage(TYPE_DELETE, data)
	if err != nil {
		// Add log to trace
		return
	}
	err = db.msgProcessor.Save(message)
	if err != nil {
		// Add log to trace
		return
	}
	// Add log to trace
	return
}

func (db *dbFullTextSearch) msgLoop() {
	for {
		msg := <-db.msgProcessor.GetMsgChannel()
		db.search.IndexData(msg)
	}
}
