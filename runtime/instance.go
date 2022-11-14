package runtime

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
	return nil
}

func (db *dbFullTextSearch) Update(tableName string, value interface{}) error {
	return nil
}

func (db *dbFullTextSearch) Delete(tableName string, value interface{}) error {
	return nil
}

func (db *dbFullTextSearch) Get(tableName string, condition string) error {
	return nil
}

func (db *dbFullTextSearch) SearchLite(query string, result interface{}) error {
	return nil
}

func (db *dbFullTextSearch) SetStorage(dbStorage DbStorage) {
	db.db = dbStorage
}

func (db *dbFullTextSearch) SetMsgProcessor(msgProcessor MessageProcessor) {
	db.msgProcessor = msgProcessor
}

func (db *dbFullTextSearch) SetSearchProcessor(search SearchProcessor) {
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
