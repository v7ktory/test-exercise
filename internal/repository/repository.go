package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/v7ktory/test/internal/model"
	"github.com/v7ktory/test/pkg/database/mongodb"
)

type Auth interface {
	Create(ctx context.Context, user *model.User) (uuid.UUID, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}
type Session interface {
	Create(ctx context.Context, session model.Session) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Session, error)
	Update(ctx context.Context, session model.Session) error
}

type Repository struct {
	Auth
	Session
}

func NewRepository(provider *mongodb.Provider) *Repository {
	return &Repository{
		Auth:    NewAuthRepository(provider),
		Session: NewSessionRepository(provider),
	}
}
