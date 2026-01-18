package storage

import (
	"context"
	"io"
	"net/url"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewS3)
}

type S3 interface {
	PutObject(ctx context.Context, objectName string, reader io.Reader, size int64) error
	GetObject(ctx context.Context, objectName string) (io.ReadCloser, error)
	PresignedGetObject(ctx context.Context, objectName string) (*url.URL, error)
}

type S3Impl struct {
	client *minio.Client
	config *config.Config
}

func NewS3(i do.Injector) (S3, error) {
	cfg := do.MustInvoke[*config.Config](i)
	client, err := minio.New(cfg.Storage.S3.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.Storage.S3.AccessKeyId, cfg.Storage.S3.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	return &S3Impl{
		client: client,
		config: cfg,
	}, nil
}

func (s *S3Impl) PutObject(ctx context.Context, objectName string, reader io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, s.config.Storage.S3.Bucket, objectName, reader, size, minio.PutObjectOptions{})
	return err
}

func (s *S3Impl) GetObject(ctx context.Context, objectName string) (io.ReadCloser, error) {
	return s.client.GetObject(ctx, s.config.Storage.S3.Bucket, objectName, minio.GetObjectOptions{})
}

func (s *S3Impl) PresignedGetObject(ctx context.Context, objectName string) (*url.URL, error) {
	return s.client.PresignedGetObject(ctx, s.config.Storage.S3.Bucket, objectName, s.config.Storage.S3.UrlExpiration, nil)
}
