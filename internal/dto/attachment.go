package dto

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/entity"
	storageS3 "github.com/anonychun/bibit/internal/storage/s3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/samber/do/v2"
)

type AttachmentBlueprint struct {
	Id       string `json:"id"`
	FileName string `json:"fileName"`
	Url      string `json:"url"`
}

func NewAttachmentBlueprint(ctx context.Context, attachment *entity.Attachment) (*AttachmentBlueprint, error) {
	if attachment == nil {
		return nil, nil
	}

	s3Storage, err := do.Invoke[storageS3.S3](bootstrap.Injector)
	if err != nil {
		return nil, err
	}

	presignResult, err := s3Storage.PresignGetObject(ctx, &s3.GetObjectInput{
		Key: aws.String(attachment.ObjectName),
	})
	if err != nil {
		return nil, err
	}

	return &AttachmentBlueprint{
		Id:       attachment.Id.String(),
		FileName: attachment.FileName,
		Url:      presignResult.URL,
	}, nil
}
