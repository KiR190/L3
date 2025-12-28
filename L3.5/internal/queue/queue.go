package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"event-booker/internal/models"

	kfg "github.com/segmentio/kafka-go"
	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
)

type QueueInterface interface {
	PublishExpiration(ctx context.Context, task models.ExpirationTask) error
	Subscribe(ctx context.Context, out chan<- kfg.Message, strategy retry.Strategy)
	Commit(ctx context.Context, msg kfg.Message) error
	Close() error
}

type Queue struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
}

func NewQueue(brokers []string, topic, groupID string) QueueInterface {
	return &Queue{
		producer: kafka.NewProducer(brokers, topic),
		consumer: kafka.NewConsumer(brokers, topic, groupID),
	}
}

func (q *Queue) PublishExpiration(ctx context.Context, task models.ExpirationTask) error {
	value, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal expiration task: %w", err)
	}

	key := []byte(task.BookingID)

	err = q.producer.Send(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	return nil
}

func (q *Queue) Subscribe(ctx context.Context, out chan<- kfg.Message, strategy retry.Strategy) {
	q.consumer.StartConsuming(ctx, out, strategy)
}

func (q *Queue) Commit(ctx context.Context, msg kfg.Message) error {
	return q.consumer.Commit(ctx, msg)
}

func (q *Queue) Close() error {
	log.Println("Closing queue connections...")

	var errs []error

	if err := q.producer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close producer: %w", err))
	}

	if err := q.consumer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close consumer: %w", err))
	}

	if len(errs) > 0 {
		log.Printf("Queue closed with errors: %v", errs)
		return fmt.Errorf("errors closing queue: %v", errs)
	}

	log.Println("Queue closed successfully.")
	return nil
}
