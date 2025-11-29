package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"

	. "image-processor/internal/models"
	"image-processor/internal/queue"
	"image-processor/internal/storage"

	"github.com/google/uuid"
)

type ImageService interface {
	UploadAndEnqueue(ctx context.Context, file *multipart.FileHeader, taskType string, params TaskParams) (string, error)
	GetImage(ctx context.Context, id string) ([]byte, error)
	DeleteImage(ctx context.Context, id string) error
	GetStatus(ctx context.Context, imageID string) (*Task, error)
}

type Service struct {
	storage storage.Storage
	queue   queue.Queue
}

func NewImageService(storage storage.Storage, queue queue.Queue) ImageService {
	return &Service{
		storage: storage,
		queue:   queue,
	}
}

// UploadAndEnqueue загружает файл и отправляет задачу на обработку
func (s *Service) UploadAndEnqueue(ctx context.Context, file *multipart.FileHeader, taskType string, params TaskParams) (string, error) {
	if file == nil {
		return "", errors.New("file is nil")
	}

	// Открываем файл
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Читаем весь файл в память
	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	// Генерируем уникальный ID
	id := uuid.New().String()

	// достаем расширение
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// Определяем контент-тип по расширению
	contentType := detectContentType(file.Filename)

	// Имя объекта: id + расширение
	objectName := id + ext
	originalPath := "original/" + objectName

	// Сохраняем в хранилище
	err = s.storage.Save(ctx, originalPath, data, contentType)
	if err != nil {
		return "", err
	}

	// формируем задачу
	task := Task{
		TaskID:      uuid.New().String(),
		ImageID:     id,
		Type:        taskType,
		Params:      params,
		Extension:   ext,
		ContentType: contentType,
		Status:      TaskStatusProcessing,
	}

	taskBytes, err := json.Marshal(task)
	if err != nil {
		return "", err
	}

	// отправка задачи в Kafka
	err = s.queue.Publish(ctx, []byte(task.ImageID), taskBytes)
	if err != nil {
		return "", err
	}

	// Сохраняем статус
	if err := s.saveTask(ctx, task); err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) saveTask(ctx context.Context, task Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return s.storage.Save(ctx,
		"status/"+task.ImageID+".json",
		data,
		"application/json",
	)
}

func (s *Service) GetStatus(ctx context.Context, imageID string) (*Task, error) {
	path := "status/" + imageID + ".json"

	data, err := s.storage.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *Service) GetImage(ctx context.Context, id string) ([]byte, error) {
	extensions := []string{".jpg", ".jpeg", ".png", ".gif"}

	for _, ext := range extensions {
		data, err := s.storage.Get(ctx, id+ext)
		if err == nil {
			return data, nil
		}
	}

	return nil, errors.New("image not found")
}

func (s *Service) DeleteImage(ctx context.Context, id string) error {
	extensions := []string{".jpg", ".jpeg", ".png", ".gif"}

	var lastErr error

	for _, ext := range extensions {
		err := s.storage.Delete(ctx, id+ext)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return lastErr
}

// detectContentType определяет MIME-type по расширению
func detectContentType(filename string) string {
	ext := filepath.Ext(filename)
	m := mime.TypeByExtension(ext)
	if m == "" {
		return "application/octet-stream"
	}
	return m
}
