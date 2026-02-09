package current

import (
	"context"

	"github.com/anonychun/bibit/internal/entity"
	"github.com/uptrace/bun"
)

type key int

const (
	txKey key = iota
	adminKey
	userKey
)

func Tx(ctx context.Context) *bun.Tx {
	tx, _ := ctx.Value(txKey).(*bun.Tx)
	return tx
}

func SetTx(ctx context.Context, tx *bun.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func Admin(ctx context.Context) *entity.Admin {
	admin, _ := ctx.Value(adminKey).(*entity.Admin)
	return admin
}

func SetAdmin(ctx context.Context, admin *entity.Admin) context.Context {
	return context.WithValue(ctx, adminKey, admin)
}

func User(ctx context.Context) *entity.User {
	user, _ := ctx.Value(userKey).(*entity.User)
	return user
}

func SetUser(ctx context.Context, user *entity.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}
