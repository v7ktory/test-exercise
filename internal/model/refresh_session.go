package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshSession struct {
	ID           uuid.UUID    `json:"id" bson:"_id"`
	UserID       uuid.UUID    `json:"user_id" bson:"user_id"`
	RefreshToken RefreshToken `json:"refresh_token" bson:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at" bson:"expires_at"`
}
