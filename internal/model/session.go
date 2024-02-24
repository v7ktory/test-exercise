package model

import (
	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID    `json:"id" bson:"_id"`
	RefreshToken RefreshToken `json:"refresh_token" bson:"refresh_token"`
}
