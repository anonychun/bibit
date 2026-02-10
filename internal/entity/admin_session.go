package entity

import (
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type AdminSession struct {
	Base

	AdminId   uuid.UUID
	Admin     *Admin `bun:"rel:belongs-to,join:admin_id=id"`
	Token     string
	IpAddress string
	UserAgent string
}

func (as *AdminSession) GenerateToken() {
	as.Token = ulid.Make().String()
}
