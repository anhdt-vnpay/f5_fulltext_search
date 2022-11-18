package message_processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	redis_connector "github.com/anhdt-vnpay/f5_fulltext_search/redis"
)

var (
	wg sync.WaitGroup
)

type messageProcessor struct {
	redisChannel   string
	redisConnector redis_connector.RedisConnector
}

func NewMessageProcessor(redisChannel string, redisConnector redis_connector.RedisConnector) *messageProcessor {
	return &messageProcessor{
		redisChannel:   redisChannel,
		redisConnector: redisConnector,
	}
}

func (r *messageProcessor) Save(data []byte) error {
	fmt.Printf("Publish message: %s > %v\n", r.redisChannel, string(data))
	r.redisConnector.Publish(r.redisChannel, data)
	return nil
}

func (r *messageProcessor) Listen() error {
	subscriber := r.redisConnector.GetClient().Subscribe(context.Background(), r.redisChannel)

	for {
		msg, err := subscriber.ReceiveMessage(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Println("Receive: ", msg.Payload)
		time.Sleep(time.Second)
	}
}
