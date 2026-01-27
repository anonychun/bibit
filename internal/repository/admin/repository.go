package admin

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	dbSql "github.com/anonychun/bibit/internal/db/sql"
	"github.com/anonychun/bibit/internal/entity"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewRepository)
}

type IRepository interface {
	FindById(ctx context.Context, id string) (*entity.Admin, error)
	FindByEmailAddress(ctx context.Context, emailAddress string) (*entity.Admin, error)
}

type Repository struct {
	sqlDB dbSql.IDB
}

var _ IRepository = (*Repository)(nil)

func NewRepository(i do.Injector) (*Repository, error) {
	return &Repository{
		sqlDB: do.MustInvoke[*dbSql.PostgresDB](i),
	}, nil
}

func (r *Repository) FindById(ctx context.Context, id string) (*entity.Admin, error) {
	admin := &entity.Admin{}
	err := r.sqlDB.DB(ctx).First(admin, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return admin, nil
}

func (r *Repository) FindByEmailAddress(ctx context.Context, emailAddress string) (*entity.Admin, error) {
	admin := &entity.Admin{}
	err := r.sqlDB.DB(ctx).First(admin, "email_address = ?", emailAddress).Error
	if err != nil {
		return nil, err
	}

	return admin, nil
}
