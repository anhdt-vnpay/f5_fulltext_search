package message_processor

import (
	"context"
	"fmt"
	"reflect"

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
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			fmt.Printf("Receive message error: %s\n", err.Error())
			subscriber = r.redisConnector.GetClient().Subscribe(ctx, r.redisChannel)
		}
		fmt.Printf("Receive message: %v and %s\n", msg.Payload, reflect.TypeOf(msg.Payload))
		r.chMsg <- []byte(msg.Payload)
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
