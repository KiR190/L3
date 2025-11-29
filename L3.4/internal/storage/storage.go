package storage

import (
	"bytes"
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage interface {
	Save(ctx context.Context, objName string, data []byte, contentType string) error
	Get(ctx context.Context, objName string) ([]byte, error)
	Delete(ctx context.Context, objectName string) error
}

type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucket string, useSSL bool) (Storage, error) {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Создание бакета если нет
	ctx := context.Background()
	err = cli.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := cli.BucketExists(ctx, bucket)
		if errBucketExists == nil && exists {
			// ok
		} else {
			return nil, err
		}
	}

	return &MinioStorage{
		client:     cli,
		bucketName: bucket,
	}, nil
}

// SaveImage сохраняет файл в S3
func (s *MinioStorage) Save(ctx context.Context, objectName string, data []byte, contentType string) error {
	reader := bytes.NewReader(data)

	_, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})

	return err
}

// GetImage получает файл
func (s *MinioStorage) Get(ctx context.Context, objectName string) ([]byte, error) {
	obj, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// DeleteImage удаляет файл
func (s *MinioStorage) Delete(ctx context.Context, objectName string) error {
	return s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
}
