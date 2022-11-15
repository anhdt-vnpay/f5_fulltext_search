package helper

import (
	"fmt"
	"strings"

	m "github.com/anhdt-vnpay/f5_fulltext_search/model"
)

func ModelByName(tableName string) (interface{}, error) {
	tableName = strings.ToLower(tableName)
	tableName = strings.Trim(tableName, " ")
	switch tableName {
	case "user":
		return m.User{}, nil
	default:
		return nil, fmt.Errorf("invalid table name")
	}
}
