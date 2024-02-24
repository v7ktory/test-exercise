package model

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	ID     uuid.UUID `json:"id" bson:"_id"`
	UserID uuid.UUID `json:"user_id" bson:"user_id"`
	Token  string    `json:"token" bson:"token"`
}

type RefreshToken struct {
	ID            uuid.UUID `json:"id" bson:"_id"`
	UserID        uuid.UUID `json:"user_id" bson:"user_id"`
	AccessTokenID uuid.UUID `json:"access_token_id" bson:"access_token_id"`
	Token         string    `json:"token" bson:"token"`
	ExpiresAt     time.Time `json:"expires_at" bson:"expires_at"`
}
