package search_processor

import (
	"context"
	"fmt"
	"time"

	redis_connector "github.com/anhdt-vnpay/f5_fulltext_search/redis"
)

type searchProcessor struct {
	ch             chan string
	redisChannel   string
	redisConnector redis_connector.RedisConnector
}

func NewSearchProcessor(ch chan string, redisChannel string, redisConnector redis_connector.RedisConnector) *searchProcessor {
	return &searchProcessor{
		ch:             ch,
		redisChannel:   redisChannel,
		redisConnector: redisConnector,
	}
}

func (r *searchProcessor) PassMessage() error {
	for {
		select {
		case msg := <-r.ch:
			fmt.Println("Pass message: ", msg)
			r.IndexData([]byte(msg))
		case <-time.After(10 * time.Second):
			fmt.Println("In pass message loop")
		}
	}
}

func (r *searchProcessor) Subscribe() error {
	subscriber := r.redisConnector.GetClient().Subscribe(context.Background(), r.redisChannel)
	for {
		msg, err := subscriber.ReceiveMessage(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Println("Receive: ", msg.Payload)
		r.ch <- msg.Payload
		time.Sleep(time.Second)
	}
}

func (r *searchProcessor) IndexData(data []byte) error {
	fmt.Println("Index data: ", string(data))
	return nil
}

func (r *searchProcessor) SearchLite(query string) ([]byte, error) {
	return nil, nil
}
