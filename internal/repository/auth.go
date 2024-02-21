package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/v7ktory/test/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type AuthRepository struct {
	db *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) *AuthRepository {
	return &AuthRepository{
		db: db.Collection("users"),
	}
}

func (r *AuthRepository) Create(ctx context.Context, user *model.User) (string, error) {
	_, err := r.db.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	return user.UUID.String(), nil
}

func (r *AuthRepository) GetByCredentials(ctx context.Context, email, password string) (*model.User, error) {
	var user model.User
	if err := r.db.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user); err != nil {
		return &model.User{}, ErrUserNotFound
	}
	return &user, nil
}

func (r *AuthRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (model.User, error) {
	return model.User{}, nil
}

func (r *AuthRepository) SetSession(ctx context.Context, userID uuid.UUID, session model.Session) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"session": session}})

	return err
}
