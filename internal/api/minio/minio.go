package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"rest_api/internal/api/config"
)

type Service struct {
	client *minio.Client
}

func (s *Service) GetObject(ctx context.Context, id string) (*minio.Object, error) {
	object, err := s.client.GetObject(ctx, config.BucketName, id, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (s *Service) PutObject(ctx context.Context, id string) error {
	return nil
}
