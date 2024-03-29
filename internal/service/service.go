package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/v7ktory/test/internal/model"
	"github.com/v7ktory/test/internal/repository"
	"github.com/v7ktory/test/pkg/hash"
	"github.com/v7ktory/test/pkg/jwt"
)

type Auth interface {
	SignUp(ctx context.Context, user *model.User) (uuid.UUID, error)
	Login(ctx context.Context, userID uuid.UUID, email, password string) (*model.AccessToken, *model.RefreshToken, error)
	Refresh(ctx context.Context, userID uuid.UUID, accessTokenBearer, refreshTokenCookie string) (*model.AccessToken, *model.RefreshToken, error)
}

type Service struct {
	Auth
}

func NewService(repo repository.Repository, hash hash.Hasher, jwt jwt.JWT, log *slog.Logger, accessTokenTTL, refreshTokenTTL time.Duration) *Service {
	return &Service{
		Auth: NewAuthService(repo, hash, jwt, log, accessTokenTTL, refreshTokenTTL),
	}
}
