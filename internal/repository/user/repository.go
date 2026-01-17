package user

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
	FindById(ctx context.Context, id string) (*entity.User, error)
	FindByEmailAddress(ctx context.Context, emailAddress string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	ExistsByEmailAddress(ctx context.Context, emailAddress string) (bool, error)
}

type RepositoryImpl struct {
	sql db.Sql
}

func NewRepository(i do.Injector) (Repository, error) {
	return &RepositoryImpl{
		sql: do.MustInvoke[db.Sql](i),
	}, nil
}

func (r RepositoryImpl) FindById(ctx context.Context, id string) (*entity.User, error) {
	user := &entity.User{}
	err := r.sql.DB(ctx).First(user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r RepositoryImpl) FindByEmailAddress(ctx context.Context, emailAddress string) (*entity.User, error) {
	user := &entity.User{}
	err := r.sql.DB(ctx).First(user, "email_address = ?", emailAddress).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r RepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	return r.sql.DB(ctx).Create(user).Error
}

func (r RepositoryImpl) ExistsByEmailAddress(ctx context.Context, emailAddress string) (bool, error) {
	var exists bool
	err := r.sql.DB(ctx).Raw("SELECT 1 FROM users WHERE email_address = ?", emailAddress).Scan(&exists).Error
	return exists, err
}
