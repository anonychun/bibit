package admin

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/db"
	"github.com/anonychun/bibit/internal/entity"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewRepository)
}

type Repository interface {
	FindById(ctx context.Context, id string) (*entity.Admin, error)
	FindByEmailAddress(ctx context.Context, emailAddress string) (*entity.Admin, error)
}

type RepositoryImpl struct {
	sql db.Sql
}

func NewRepository(i do.Injector) (Repository, error) {
	return &RepositoryImpl{
		sql: do.MustInvoke[db.Sql](i),
	}, nil
}

func (r *RepositoryImpl) FindById(ctx context.Context, id string) (*entity.Admin, error) {
	admin := &entity.Admin{}
	err := r.sql.DB(ctx).First(admin, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return admin, nil
}

func (r *RepositoryImpl) FindByEmailAddress(ctx context.Context, emailAddress string) (*entity.Admin, error) {
	admin := &entity.Admin{}
	err := r.sql.DB(ctx).First(admin, "email_address = ?", emailAddress).Error
	if err != nil {
		return nil, err
	}

	return admin, nil
}
