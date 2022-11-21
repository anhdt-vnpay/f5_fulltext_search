package helper

import (
	"fmt"
	"strings"

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

func ParseMessage(message string) (string, string, string, error) {
	strArr := strings.Split(message, "/")
	if len(strArr) < 3 {
		return "", "", "", fmt.Errorf("parse invalid message")
	}
	return strArr[0], strArr[1], strArr[2], nil
}

func GetDataID(data string) string {
	data = strings.Trim(data, "{")
	data = strings.Trim(data, "}")

	dataArr := strings.Split(data, ",")
	if len(dataArr) < 1 {
		return "something's wrong"
	}

	idPart := strings.Split(dataArr[0], ":")
	if len(idPart) < 2 {
		return "something's wrong too"
	}

	return fmt.Sprint(idPart[1])
}
