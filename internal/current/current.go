package current

import (
	"context"

	"github.com/anonychun/bibit/internal/entity"
	"github.com/uptrace/bun"
)

type key int

const (
	txKey key = iota
	userKey
)

func Tx(ctx context.Context) *bun.Tx {
	tx, _ := ctx.Value(txKey).(*bun.Tx)
	return tx
}

func SetTx(ctx context.Context, tx *bun.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func User(ctx context.Context) *entity.User {
	user, _ := ctx.Value(userKey).(*entity.User)
	return user
}

func SetUser(ctx context.Context, user *entity.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}
