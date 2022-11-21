package search_processor

import (
	"fmt"
	"time"

	redis_connector "github.com/anhdt-vnpay/f5_fulltext_search/redis"
	handler "github.com/anhdt-vnpay/f5_fulltext_search/search_processor/handler"
	helper "github.com/anhdt-vnpay/f5_fulltext_search/search_processor/helper"
	"github.com/olivere/elastic/v7"
)

type searchProcessor struct {
	ch             chan []byte
	redisChannel   string
	redisConnector redis_connector.RedisConnector
	es             *elastic.Client
}

func NewSearchProcessor(ch chan []byte, redisChannel string, redisConnector redis_connector.RedisConnector) *searchProcessor {
	es, err := helper.GetESClient()
	if err != nil {
		fmt.Printf("error initialize Elastic client: %s\n", err.Error())
		return nil
	}

	searchProcessor := &searchProcessor{
		ch:             ch,
		redisChannel:   redisChannel,
		redisConnector: redisConnector,
		es:             es,
	}
	searchProcessor.PassMessage()
	return searchProcessor
}

func (r *searchProcessor) PassMessage() error {
	pass := func() {
		for {
			select {
			case msg := <-r.ch:
				fmt.Println("Pass message: ", string(msg))
				r.IndexData(msg)
			case <-time.After(10 * time.Second):
				fmt.Println("In pass message loop")
			}
		}
	}
	go func() {
		pass()
	}()
	return nil
}

func (r *searchProcessor) IndexData(data []byte) error {
	tipe, index, body, err := helper.ParseMessage(data)
	if err != nil {
		return err
	}
	id, err := helper.GetDataID(body)
	if err != nil {
		return err
	}

	switch tipe {
	case "insert":
		fmt.Println("ES insert >>")
		if err := handler.Insert(r.es, index, id, body); err != nil {
			fmt.Printf("es insert error: %s\n", err.Error())
			return err
		}
	case "update":
		fmt.Println("ES update >>")
		if err := handler.Update(r.es, index, id, body); err != nil {
			fmt.Printf("es update error: %s\n", err.Error())
			return err
		}
	case "delete":
		fmt.Println("ES delete >>")
		if err := handler.Delete(r.es, index, id, body); err != nil {
			fmt.Printf("es delete error: %s\n", err.Error())
			return err
		}
	default:
		fmt.Printf("abnormal message type\n")
	}

	return nil
}

func (r *searchProcessor) SearchLite(query string) ([]byte, error) {
	rs, err := handler.Search(r.es, query)
	if err != nil {
		fmt.Printf("search lite error: %s\n", err.Error())
		return nil, err
	}
	return rs, nil
}
