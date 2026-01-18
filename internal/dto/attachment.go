package dto

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/entity"
	"github.com/anonychun/bibit/internal/storage"
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

	s3, err := do.Invoke[storage.S3](bootstrap.Injector)
	if err != nil {
		return nil, err
	}

	u, err := s3.PresignedGetObject(ctx, attachment.ObjectName)
	if err != nil {
		return nil, err
	}

	return &AttachmentBlueprint{
		Id:       attachment.Id.String(),
		FileName: attachment.FileName,
		Url:      u.String(),
	}, nil
}
