package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/v7ktory/test/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Auth interface {
	Create(ctx context.Context, user *model.User) (uuid.UUID, error)
	GetByCredentials(ctx context.Context, email, password string) (*model.User, error)
	SetSession(ctx context.Context, userID uuid.UUID, session model.Session) error
}

type Repository struct {
	Auth
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Auth: NewAuthRepository(db),
	}
}
