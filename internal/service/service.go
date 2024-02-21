package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/v7ktory/test/internal/model"
	"github.com/v7ktory/test/internal/repository"
	"github.com/v7ktory/test/pkg/hash"
	"github.com/v7ktory/test/pkg/jwt"
)

type Auth interface {
	SignUp(ctx context.Context, user *model.User) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type Service struct {
	Auth
}

func NewService(repo repository.AuthRepository, hash hash.Hasher, jwt jwt.JWT, log *slog.Logger, ttl time.Duration) *Service {
	return &Service{
		Auth: NewAuthService(repo, hash, jwt, log, ttl),
	}
}
