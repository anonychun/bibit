package user

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
	FindById(ctx context.Context, id string) (*entity.User, error)
	FindByEmailAddress(ctx context.Context, emailAddress string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	ExistsByEmailAddress(ctx context.Context, emailAddress string) (bool, error)
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

func (r *Repository) FindById(ctx context.Context, id string) (*entity.User, error) {
	user := &entity.User{}
	err := r.sqlDB.DB(ctx).First(user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) FindByEmailAddress(ctx context.Context, emailAddress string) (*entity.User, error) {
	user := &entity.User{}
	err := r.sqlDB.DB(ctx).First(user, "email_address = ?", emailAddress).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) Create(ctx context.Context, user *entity.User) error {
	return r.sqlDB.DB(ctx).Create(user).Error
}

func (r *Repository) ExistsByEmailAddress(ctx context.Context, emailAddress string) (bool, error) {
	var exists bool
	err := r.sqlDB.DB(ctx).Raw("SELECT 1 FROM users WHERE email_address = ?", emailAddress).Scan(&exists).Error
	return exists, err
}
