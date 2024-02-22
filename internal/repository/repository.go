package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/v7ktory/test/internal/model"
	"github.com/v7ktory/test/pkg/database/mongodb"
)

type Auth interface {
	Create(ctx context.Context, user *model.User) (uuid.UUID, error)
	GetByCredentials(ctx context.Context, email, password string) (*model.User, error)
	SetSession(ctx context.Context, userID uuid.UUID, session model.RefreshSession) error
}

type Repository struct {
	Auth
}

func NewRepository(provider *mongodb.Provider) *Repository {
	return &Repository{
		Auth: NewAuthRepository(provider),
	}
}
