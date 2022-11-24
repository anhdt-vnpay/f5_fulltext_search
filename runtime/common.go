package runtime

import (
	"encoding/json"
	"fmt"
	"strings"

	model "github.com/anhdt-vnpay/f5_fulltext_search/model"
)

func createMessage(tipe string, tableName string, data interface{}) ([]byte, error) {
	tipe = strings.Trim(tipe, " ")
	tipe = strings.ToLower(tipe)

	tableName = strings.Trim(tableName, " ")
	tableName = strings.ToLower(tableName)

	m := model.Message{
		Tipe:      tipe,
		TableName: tableName,
		Data:      data,
	}

	byteData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return byteData, nil
}

func parseSearchResult(data []byte) ([]map[string]interface{}, error) {
	var rs []map[string]interface{}
	fmt.Println("DEBUG data: ", data)
	err := json.Unmarshal(data, &rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
