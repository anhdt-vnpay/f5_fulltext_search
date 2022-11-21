package runtime

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TODO: need to refactor => create message object
func createMessage(tipe string, tableName string, data interface{}) ([]byte, error) {
	tipe = strings.Trim(tipe, " ")
	tipe = strings.ToLower(tipe)

	tableName = strings.Trim(tableName, " ")
	tableName = strings.ToLower(tableName)

	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	m := fmt.Sprintf("%s/%s/%s", tipe, tableName, string(byteData))
	fmt.Println("message: ", m)
	return []byte(m), nil
}

func parseSearchResult(data []byte) ([]map[string]interface{}, error) {
	var rs []map[string]interface{}
	err := json.Unmarshal(data, &rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
