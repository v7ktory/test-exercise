package model

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	Token  string
	ID     uuid.UUID
	UserID uuid.UUID
}

type RefreshToken struct {
	Token     string
	ID        uuid.UUID
	UserID    uuid.UUID
	ExpiresAt time.Time
}
