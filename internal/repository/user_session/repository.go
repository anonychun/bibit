package user_session

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
	FindByToken(ctx context.Context, token string) (*entity.UserSession, error)
	Create(ctx context.Context, userSession *entity.UserSession) error
	DeleteByToken(ctx context.Context, token string) error
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

func (r *Repository) FindByToken(ctx context.Context, token string) (*entity.UserSession, error) {
	userSession := &entity.UserSession{}
	err := r.sqlDB.DB(ctx).NewSelect().Model(userSession).Where("token = ?", token).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return userSession, nil
}

func (r *Repository) Create(ctx context.Context, userSession *entity.UserSession) error {
	_, err := r.sqlDB.DB(ctx).NewInsert().Model(userSession).Exec(ctx)
	return err
}

func (r *Repository) DeleteByToken(ctx context.Context, token string) error {
	_, err := r.sqlDB.DB(ctx).NewDelete().Model(&entity.UserSession{}).Where("token = ?", token).Exec(ctx)
	return err
}
