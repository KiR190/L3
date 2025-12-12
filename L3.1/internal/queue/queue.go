package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"delayed-notifier/internal/models"

	amqp091 "github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
)

type Queue struct {
	conn      *rabbitmq.Connection
	channel   *rabbitmq.Channel
	publisher *rabbitmq.Publisher
	consumer  *rabbitmq.Consumer
	stopChan  chan struct{}
	handler   func(context.Context, models.Notification) error
	wg        sync.WaitGroup
}

func NewQueue(url string) (*Queue, error) {
	conn, err := rabbitmq.Connect(url, 3, 5*time.Second)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	q := &Queue{
		conn:     conn,
		channel:  ch,
		stopChan: make(chan struct{}),
	}

	// Настраиваем очереди
	if err := q.SetupQueues(); err != nil {
		_ = conn.Close()
		return nil, err
	}

	// Создаем publisher и consumer
	q.publisher = rabbitmq.NewPublisher(ch, "")
	q.consumer = rabbitmq.NewConsumer(ch, rabbitmq.NewConsumerConfig("notifications.main"))

	return q, nil
}

// Consume читает сообщения из очереди.
func (q *Queue) Consume(msgChan chan []byte) error {
	return q.consumer.Consume(msgChan)
}

func (q *Queue) SetupQueues() error {
	// Создаём delayed exchange
	args := amqp091.Table{
		"x-delayed-type": "direct",
	}
	if err := q.channel.ExchangeDeclare(
		"notifications.delayed", // имя exchange
		"x-delayed-message",     // тип
		true,                    // durable
		false,                   // autoDelete
		false,                   // internal
		false,                   // noWait
		args,                    // arguments
	); err != nil {
		return err
	}

	// Создаём DLQ
	if _, err := q.channel.QueueDeclare(
		"notifications.dlq",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	); err != nil {
		return err
	}

	// Создаём main очередь с DLQ
	mainQueueArgs := amqp091.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": "notifications.dlq",
	}
	if _, err := q.channel.QueueDeclare(
		"notifications.main",
		true,
		false,
		false,
		false,
		mainQueueArgs,
	); err != nil {
		return err
	}

	// Привязываем main очередь к delayed exchange
	if err := q.channel.QueueBind(
		"notifications.main",    // queue
		"notifications.main",    // routing key
		"notifications.delayed", // exchange
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}

// Публикация
func (q *Queue) Publish(body []byte, sendAt time.Time) error {
	delay := time.Until(sendAt)
	if delay < 0 {
		delay = 0
	}

	err := q.channel.Publish(
		"notifications.delayed", // exchange
		"notifications.main",    // routing key
		false,                   // mandatory
		false,                   // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
			Headers: amqp091.Table{
				"x-delay": int64(delay.Milliseconds()),
			},
		},
	)

	if err != nil {
		log.Printf("Publish error: %v", err)
	} else {
		log.Printf("Message published to exchange with delay %dms", delay.Milliseconds())
	}

	return err
}

// Запуск DLQ Consumer
func (q *Queue) StartDLQConsumer() {
	msgs, err := q.channel.Consume(
		"notifications.dlq",
		"",
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		log.Println("DLQ consume error:", err)
		return
	}

	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		for {
			select {
			case <-q.stopChan:
				log.Println("DLQ consumer stopping...")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("DLQ consumer channel closed")
					return
				}
				log.Printf("DLQ сообщение: %s\n", string(msg.Body))
			}
		}
	}()
}

// Запуск Main Consumer
func (q *Queue) StartMainConsumer() {
	msgs, err := q.channel.Consume(
		"notifications.main", // main очередь
		"",                   // consumer tag
		true,                 // autoAck
		false,                // exclusive
		false,                // noLocal
		false,                // noWait
		nil,                  // args
	)
	if err != nil {
		log.Println("Main consume error:", err)
		return
	}

	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		for {
			select {
			case <-q.stopChan:
				log.Println("Main consumer stopping...")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("Main consumer channel closed")
					return
				}
				var notification models.Notification
				if err := json.Unmarshal(msg.Body, &notification); err != nil {
					log.Println("Failed to parse notification:", err)
					continue
				}

				if q.handler != nil {
					if err := q.handler(context.Background(), notification); err != nil {
						log.Println("handler failed:", err)
					}
				} else {
					log.Println("no handler set for queue")
				}
			}
		}
	}()
}

func (q *Queue) SetHandler(handler func(context.Context, models.Notification) error) {
	q.handler = handler
}

// Close останавливает всех консьюмеров и закрывает соединения
func (q *Queue) Close() error {
	log.Println("Closing queue consumers...")

	close(q.stopChan)
	q.wg.Wait()

	if err := q.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := q.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	log.Println("Queue closed.")
	return nil
}
