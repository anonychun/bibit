package user_session

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
	FindByToken(ctx context.Context, token string) (*entity.UserSession, error)
	Create(ctx context.Context, userSession *entity.UserSession) error
	DeleteByToken(ctx context.Context, token string) error
}

type RepositoryImpl struct {
	sql db.Sql
}

func NewRepository(i do.Injector) (Repository, error) {
	return &RepositoryImpl{
		sql: do.MustInvoke[db.Sql](i),
	}, nil
}

func (r *RepositoryImpl) FindByToken(ctx context.Context, token string) (*entity.UserSession, error) {
	userSession := &entity.UserSession{}
	err := r.sql.DB(ctx).First(userSession, "token = ?", token).Error
	if err != nil {
		return nil, err
	}

	return userSession, nil
}

func (r *RepositoryImpl) Create(ctx context.Context, userSession *entity.UserSession) error {
	return r.sql.DB(ctx).Create(userSession).Error
}

func (r *RepositoryImpl) DeleteByToken(ctx context.Context, token string) error {
	return r.sql.DB(ctx).Delete(&entity.UserSession{}, "token = ?", token).Error
}
