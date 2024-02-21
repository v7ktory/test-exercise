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

type AuthService struct {
	repo repository.AuthRepository
	jwt  jwt.JWT
	hash hash.Hasher
	log  *slog.Logger
	ttl  time.Duration
}

func NewAuthService(repo repository.AuthRepository, hash hash.Hasher, jwt jwt.JWT, log *slog.Logger, ttl time.Duration) *AuthService {
	return &AuthService{
		repo: repo,
		jwt:  jwt,
		hash: hash,
		log:  log,
		ttl:  ttl,
	}
}

func (s *AuthService) SignUp(ctx context.Context, user *model.User) (string, error) {
	if err := user.Validate(); err != nil {
		s.log.Error("failed to validate user", err)
		return "", err
	}

	hashedPassword, err := s.hash.HashPassword(user.Password)
	if err != nil {
		s.log.Error("failed to hash password", err)
		return "", err
	}

	u := model.User{
		UUID:         uuid.New(),
		Name:         user.Name,
		Email:        user.Email,
		Password:     hashedPassword,
		RegisteredAt: time.Now(),
	}

	userID, err := s.repo.Create(ctx, &u)
	if err != nil {
		s.log.Error("failed to create user", err)
		return "", err
	}

	access, refresh, err := s.jwt.GenerateTokenPair(userID, s.ttl)
	if err != nil {
		s.log.Error("failed to generate token pair", err)
		return "", err
	}

	err = s.repo.SetSession(ctx, u.UUID, model.Session{
		RefreshToken: refresh,
		ExpiresAt:    time.Now().Add(s.ttl),
	})
	if err != nil {
		s.log.Error("failed to set session", err)
		return "", err
	}

	s.log.Info("user created successfully")
	return access, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	return "", nil
}
