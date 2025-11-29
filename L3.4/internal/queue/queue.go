package queue

import (
	"context"

	kfg "github.com/segmentio/kafka-go"

	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
)

type Queue struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
	out      chan kfg.Message
}

func NewQueue(brokers []string, topic, groupID string) Queue {
	return Queue{
		producer: kafka.NewProducer(brokers, topic),
		consumer: kafka.NewConsumer(brokers, topic, groupID),
		out:      make(chan kfg.Message),
	}
}

func (q *Queue) Publish(ctx context.Context, key, value []byte) error {
	return q.producer.Send(ctx, key, value)
}

// Subscribe запускает consumer и возвращает канал
func (q *Queue) Subscribe(ctx context.Context, strategy retry.Strategy) <-chan kfg.Message {
	go q.consumer.StartConsuming(ctx, q.out, strategy)
	return q.out
}

// Commit подтверждает обработку
func (q *Queue) Commit(ctx context.Context, msg kfg.Message) error {
	return q.consumer.Commit(ctx, msg)
}
