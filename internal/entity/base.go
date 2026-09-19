package entity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Base struct {
	Id        uuid.UUID `bun:"id,pk,type:uuid,default:uuidv7()"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:now()"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:now()"`
}

func (b *Base) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	b.UpdatedAt = time.Now()
	return nil
}
