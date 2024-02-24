package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/v7ktory/test/internal/model"
	"github.com/v7ktory/test/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type AuthRepository struct {
	provider *mongodb.Provider
}

func NewAuthRepository(provider *mongodb.Provider) *AuthRepository {
	return &AuthRepository{
		provider: provider,
	}
}

func (r *AuthRepository) Create(ctx context.Context, user *model.User) (uuid.UUID, error) {
	collection := r.provider.GetCollection("users")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(r.provider.QueryTimeout))
	defer cancel()

	if err := collection.FindOne(ctx, bson.M{"email": user.Email}).Err(); err == nil {
		return uuid.Nil, ErrUserExists
	}

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}
	return user.UUID, nil
}

func (r *AuthRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	collection := r.provider.GetCollection("users")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(r.provider.QueryTimeout))
	defer cancel()

	var user model.User
	if err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}
