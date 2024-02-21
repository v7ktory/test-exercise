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
	repo            repository.Repository
	jwt             jwt.JWT
	hash            hash.Hasher
	log             *slog.Logger
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthService(repo repository.Repository, hash hash.Hasher, jwt jwt.JWT, log *slog.Logger, accessTokenTTL, refreshTokenTTL time.Duration) *AuthService {
	return &AuthService{
		repo:            repo,
		jwt:             jwt,
		hash:            hash,
		log:             log,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AuthService) SignUp(ctx context.Context, user *model.User) (string, string, error) {
	if err := user.Validate(); err != nil {
		s.log.Error("failed to validate user", err)
		return "", "", err
	}

	hashedPassword, err := s.hash.Hash(user.Password)
	if err != nil {
		s.log.Error("failed to hash password", err)
		return "", "", err
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
		return "", "", err
	}

	access, refresh, err := s.jwt.GenerateTokenPair(userID, s.accessTokenTTL, s.refreshTokenTTL)
	if err != nil {
		s.log.Error("failed to generate token pair", err)
		return "", "", err
	}

	hashedRefresh, err := s.hash.Hash(refresh.Token)
	if err != nil {
		s.log.Error("failed to hash refresh token", err)
		return "", "", err
	}

	err = s.repo.SetSession(ctx, u.UUID, model.Session{
		ID:     uuid.New(),
		UserID: u.UUID,
		RefreshToken: model.RefreshToken{
			Token:     hashedRefresh,
			ID:        refresh.ID,
			UserID:    u.UUID,
			ExpiresAt: time.Now().Add(s.refreshTokenTTL),
		},
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
	})
	if err != nil {
		s.log.Error("failed to set session", err)
		return "", "", err
	}

	s.log.Info("user created successfully")
	return access.Token, refresh.Token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	return "", nil
}
