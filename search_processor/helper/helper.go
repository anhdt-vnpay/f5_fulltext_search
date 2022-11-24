package helper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anhdt-vnpay/f5_fulltext_search/model"
	elastic "github.com/olivere/elastic/v7"
)

var (
	ES_BASE_URL = "http://localhost:9200"
)

func GetESClient() (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL(ES_BASE_URL),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	fmt.Println("ES initialized...")

	return client, err

}

func ParseMessage(message []byte) (string, string, string, error) {
	var ori model.Message
	err := json.Unmarshal(message, &ori)
	if err != nil {
		return "", "", "", err
	}

	byteData, err := json.Marshal(ori.Data)
	if err != nil {
		return "", "", "", err
	}
	return ori.Tipe, ori.TableName, string(byteData), nil
}

func GetDataID(data string) (string, error) {
	data = strings.Trim(data, "{")
	data = strings.Trim(data, "}")

	dataArr := strings.Split(data, ",")
	if len(dataArr) < 1 {
		return "", fmt.Errorf("something's wrong")
	}

	for _, item := range dataArr {
		arr := strings.Split(item, ":")
		if len(arr) < 2 {
			return "", fmt.Errorf("something's wrong")
		}
		key := arr[0]
		value := arr[1]

		key = strings.Trim(key, "\"") // trim quotes
		key = strings.ToLower(key)

		if key == "id" {
			return fmt.Sprint(value), nil
		}
	}

	return "", fmt.Errorf("id not found")
}
