package model

type Message struct {
	Tipe      string      `json:"tipe"`
	TableName string      `json:"table_name"`
	Data      interface{} `json:"data"`
}
