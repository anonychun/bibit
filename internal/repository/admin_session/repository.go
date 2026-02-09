package admin_session

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
	FindByToken(ctx context.Context, token string) (*entity.AdminSession, error)
	Create(ctx context.Context, adminSession *entity.AdminSession) error
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

func (r *Repository) FindByToken(ctx context.Context, token string) (*entity.AdminSession, error) {
	adminSession := &entity.AdminSession{}
	err := r.sqlDB.DB(ctx).NewSelect().Model(adminSession).Where("token = ?", token).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return adminSession, nil
}

func (r *Repository) Create(ctx context.Context, adminSession *entity.AdminSession) error {
	_, err := r.sqlDB.DB(ctx).NewInsert().Model(adminSession).Exec(ctx)
	return err
}

func (r *Repository) DeleteByToken(ctx context.Context, token string) error {
	_, err := r.sqlDB.DB(ctx).NewDelete().Model(&entity.AdminSession{}).Where("token = ?", token).Exec(ctx)
	return err
}
