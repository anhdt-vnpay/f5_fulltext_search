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
	SearchLite(query string) (interface{}, error)
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
	go db.saveMsg(TYPE_INSERT, tableName, value)
	return nil
}

func (db *dbFullTextSearch) Update(tableName string, value interface{}) error {
	if err := db.db.Update(tableName, value); err != nil {
		return err
	}
	go db.saveMsg(TYPE_UPDATE, tableName, value)
	return nil
}

func (db *dbFullTextSearch) Delete(tableName string, value interface{}) error {
	if err := db.db.Delete(tableName, value); err != nil {
		return err
	}
	go db.saveMsg(TYPE_DELETE, tableName, value)
	return nil
}
func (db *dbFullTextSearch) Get(tableName string, condition string, result interface{}) error {
	fmt.Println("Get >>")
	err := db.db.Get(tableName, condition, result)
	return err
}

func (db *dbFullTextSearch) SearchLite(query string) (interface{}, error) {
	fmt.Println("Search lite >>")
	bRs, err := db.search.SearchLite(query)
	if err != nil {
		fmt.Println("Search lite error: ", err.Error())
		return nil, err
	}

	rs, err := parseSearchResult(bRs)
	if err != nil {
		fmt.Println("Parse search result error: ", err.Error())
	}

	return rs, nil
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

func (db *dbFullTextSearch) saveMsg(tipe string, tableName string, data interface{}) {
	message, err := createMessage(tipe, tableName, data)
	if err != nil {
		fmt.Printf("create message error: %s\n", err.Error())
		return
	}
	err = db.msgProcessor.Save(message)
	if err != nil {
		fmt.Printf("unknown error\n")
		return
	}
}

func (db *dbFullTextSearch) msgLoop() {
	for {
		msg := <-db.msgProcessor.GetMsgChannel()
		db.search.IndexData(msg)
	}
}
