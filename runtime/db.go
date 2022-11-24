package runtime

type DbStorage interface {
	Insert(tableName string, value interface{}) error
	Update(tableName string, value interface{}) error
	Delete(tableName string, value interface{}) error
	Get(tableName string, condition string, result interface{}) error
}
