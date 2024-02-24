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
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrAccessTokenNotFound  = errors.New("access token not found")
)

type SessionRepository struct {
	provider *mongodb.Provider
}

func NewSessionRepository(provider *mongodb.Provider) *SessionRepository {
	return &SessionRepository{
		provider: provider,
	}
}

// Создаём сессию
func (r *SessionRepository) Create(ctx context.Context, session model.Session) error {
	collection := r.provider.GetCollection("sessions")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(r.provider.QueryTimeout))
	defer cancel()

	_, err := collection.InsertOne(ctx, session)
	if err != nil {
		return err
	}
	return nil
}

// Возвращаем сессию по ID
func (r *SessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Session, error) {
	collection := r.provider.GetCollection("sessions")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(r.provider.QueryTimeout))
	defer cancel()

	var session model.Session
	err := collection.FindOne(ctx, bson.M{"refresh_token.user_id": userID}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Обновляем сессию
func (r *SessionRepository) Update(ctx context.Context, session model.Session) error {
	collection := r.provider.GetCollection("sessions")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(r.provider.QueryTimeout))
	defer cancel()

	filter := bson.M{"_id": session.ID}
	update := bson.M{"$set": bson.M{
		"refresh_token": bson.M{
			"_id":             session.RefreshToken.ID,
			"user_id":         session.RefreshToken.UserID,
			"access_token_id": session.RefreshToken.AccessTokenID,
			"token":           session.RefreshToken.Token,
			"expires_at":      session.RefreshToken.ExpiresAt,
		},
	}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
