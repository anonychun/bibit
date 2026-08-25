package health

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewUsecase)
}

type IUsecase interface {
	Up(ctx context.Context) (*UpResponse, error)
}

type Usecase struct {
}

var _ IUsecase = (*Usecase)(nil)

func NewUsecase(i do.Injector) (*Usecase, error) {
	return &Usecase{}, nil
}

func (u *Usecase) Up(ctx context.Context) (*UpResponse, error) {
	return &UpResponse{
		Status: "UP",
	}, nil
}
