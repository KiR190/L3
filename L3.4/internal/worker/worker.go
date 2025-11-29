package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/wb-go/wbf/retry"

	. "image-processor/internal/models"
	"image-processor/internal/processor"
	"image-processor/internal/queue"
	"image-processor/internal/storage"
)

type Worker struct {
	storage   storage.Storage
	processor *processor.Processor
	queue     queue.Queue
}

func NewWorker(storage storage.Storage, processor *processor.Processor, queue queue.Queue) Worker {
	return Worker{
		storage:   storage,
		processor: processor,
		queue:     queue,
	}
}

// Start запускает фонового воркера
func (w *Worker) Start() {
	fmt.Println("Worker started...")

	// Стратегия повторов
	strategy := retry.Strategy{
		Attempts: 5, // число попыток
		Delay:    1 * time.Minute,
		Backoff:  2, // экспоненциальная задержка
	}

	// запускаем чтение сообщений
	msgChan := w.queue.Subscribe(context.Background(), strategy)

	go func() {
		for msg := range msgChan {
			// декодируем задачу
			var task Task
			if err := json.Unmarshal(msg.Value, &task); err != nil {
				fmt.Println("Invalid task format:", err)
				_ = w.queue.Commit(context.Background(), msg)
				continue
			}

			fmt.Println("Processing task:", task.TaskID)

			if err := w.handleTask(task); err != nil {
				fmt.Printf("Task failed: %v\n", err)

				_ = w.queue.Commit(context.Background(), msg)
				continue
			}

			fmt.Println("Task completed:", task.TaskID)

			// подтверждаем обработку
			if err := w.queue.Commit(context.Background(), msg); err != nil {
				fmt.Println("Commit error:", err)
			}
		}
	}()
}

func (w *Worker) handleTask(task Task) error {
	ctx := context.Background()

	originalPath := "original/" + task.ImageID + task.Extension

	// скачиваем оригинальное изображение
	data, err := w.storage.Get(ctx, originalPath)
	if err != nil {
		w.markFailed(task)
		return fmt.Errorf("download error: %w", err)
	}

	// обрабатываем
	out, err := w.processor.Process(task, data)
	if err != nil {
		w.markFailed(task)
		return fmt.Errorf("processing error: %w", err)
	}

	// сохраняем
	resultPath := task.ImageID + task.Extension

	err = w.storage.Save(ctx, resultPath, out, task.ContentType)
	if err != nil {
		w.markFailed(task)
		return fmt.Errorf("upload error: %w", err)
	}

	w.markDone(task)

	return nil
}

func (w *Worker) saveTask(task Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	path := "status/" + task.ImageID + ".json"
	return w.storage.Save(context.Background(), path, data, "application/json")
}

func (w *Worker) markDone(task Task) error {
	task.Status = TaskStatusDone
	return w.saveTask(task)
}

func (w *Worker) markFailed(task Task) error {
	task.Status = TaskStatusFailed
	return w.saveTask(task)
}
