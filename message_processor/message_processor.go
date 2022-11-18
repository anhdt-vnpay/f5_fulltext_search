package message_processor

import (
	"context"
	"fmt"

	redis_connector "github.com/anhdt-vnpay/f5_fulltext_search/redis"
)

type messageProcessor struct {
	redisChannel   string
	redisConnector redis_connector.RedisConnector

	chMsg chan []byte
}

func NewMessageProcessor(redisChannel string, redisConnector redis_connector.RedisConnector) *messageProcessor {
	messageProcessor := &messageProcessor{
		chMsg:          make(chan []byte),
		redisChannel:   redisChannel,
		redisConnector: redisConnector,
	}
	messageProcessor.subscribe()
	return messageProcessor
}

func (r *messageProcessor) Save(data []byte) error {
	fmt.Printf("Publish message: %s > %v\n", r.redisChannel, string(data))
	r.redisConnector.Publish(r.redisChannel, data)
	return nil
}

func (r *messageProcessor) subscribe() {
	ctx := context.Background()
	subscriber := r.redisConnector.GetClient().Subscribe(ctx, r.redisChannel)
	receivingMsg := func() {
		msg, err := subscriber.Receive(ctx)
		if err != nil {
			subscriber = r.redisConnector.GetClient().Subscribe(ctx, r.redisChannel)
		}
		if data, ok := msg.([]byte); ok {
			// TODO: add timeout for messsage?
			r.chMsg <- data
		}
	}
	receivingLoop := func() {
		for {
			receivingMsg()
		}
	}
	go receivingLoop()
}

func (r *messageProcessor) GetMsgChannel() chan []byte {
	return r.chMsg
}
