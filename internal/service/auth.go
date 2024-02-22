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
	hash            hash.Hasher
	jwt             jwt.JWT
	log             *slog.Logger
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthService(repo repository.Repository, hash hash.Hasher, jwt jwt.JWT, log *slog.Logger, accessTokenTTL, refreshTokenTTL time.Duration) *AuthService {
	return &AuthService{
		repo:            repo,
		hash:            hash,
		jwt:             jwt,
		log:             log,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AuthService) SignUp(ctx context.Context, user *model.User) (*model.AccessToken, *model.RefreshToken, error) {
	hashedPassword, err := s.hash.Hash(user.Password)
	if err != nil {
		s.log.Error("failed to hash password:", err)
		return nil, nil, err
	}

	user.UUID = uuid.New()
	user.Password = hashedPassword

	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		s.log.Error("failed to create user:", err)
		return nil, nil, err
	}

	access, refresh, err := s.jwt.GenerateTokenPair(userID, s.accessTokenTTL, s.refreshTokenTTL)
	if err != nil {
		s.log.Error("failed to generate token pair:", err)
		return nil, nil, err
	}

	hashedRefresh, err := s.hash.Hash(refresh.Token)
	if err != nil {
		s.log.Error("failed to hash refresh token:", err)
		return nil, nil, err
	}

	refreshSession := model.RefreshSession{
		ID:     uuid.New(),
		UserID: userID,
		RefreshToken: model.RefreshToken{
			Token:     hashedRefresh,
			ID:        refresh.ID,
			UserID:    userID,
			ExpiresAt: time.Now().Add(s.refreshTokenTTL),
		},
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, userID, refreshSession)
	if err != nil {
		s.log.Error("failed to set session:", err)
		return nil, nil, err
	}

	s.log.Info("user created successfully")
	return access, refresh, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*model.AccessToken, *model.RefreshToken, error) {
	user, err := s.repo.GetByCredentials(ctx, email, password)
	if err != nil {
		s.log.Error("failed to get user by credentials:", err)
		return nil, nil, err
	}

	access, refresh, err := s.jwt.GenerateTokenPair(user.UUID, s.accessTokenTTL, s.refreshTokenTTL)
	if err != nil {
		s.log.Error("failed to generate token pair:", err)
		return nil, nil, err
	}

	hashedRefresh, err := s.hash.Hash(refresh.Token)
	if err != nil {
		s.log.Error("failed to hash refresh token:", err)
		return nil, nil, err
	}

	refreshSession := model.RefreshSession{
		ID:     uuid.New(),
		UserID: user.UUID,
		RefreshToken: model.RefreshToken{
			Token:     hashedRefresh,
			ID:        refresh.ID,
			UserID:    user.UUID,
			ExpiresAt: time.Now().Add(s.refreshTokenTTL),
		},
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, user.UUID, refreshSession)
	if err != nil {
		s.log.Error("failed to set session:", err)
		return nil, nil, err
	}

	s.log.Info("user logged in successfully")
	return access, refresh, nil
}
