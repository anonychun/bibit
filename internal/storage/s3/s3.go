package s3

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewS3)
}

type S3 interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

type S3Impl struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	config        *config.Config
}

func NewS3(i do.Injector) (S3, error) {
	ctx := context.Background()
	cfg := do.MustInvoke[*config.Config](i)

	awsCfg, err := awsConfig.LoadDefaultConfig(
		ctx,
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.Storage.S3.AccessKeyId,
			cfg.Storage.S3.SecretAccessKey,
			"",
		)),
		awsConfig.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Storage.S3.Endpoint)
	})

	return &S3Impl{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		config:        cfg,
	}, nil
}

func (s *S3Impl) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	if params.Bucket == nil {
		params.Bucket = aws.String(s.config.Storage.S3.Bucket)
	}

	return s.client.PutObject(ctx, params, optFns...)
}

func (s *S3Impl) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if params.Bucket == nil {
		params.Bucket = aws.String(s.config.Storage.S3.Bucket)
	}

	return s.client.GetObject(ctx, params, optFns...)
}

func (s *S3Impl) PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	if params.Bucket == nil {
		params.Bucket = aws.String(s.config.Storage.S3.Bucket)
	}

	return s.presignClient.PresignGetObject(ctx, params, optFns...)
}
