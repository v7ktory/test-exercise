package service

import (
	"context"
	"errors"
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

func (s *AuthService) SignUp(ctx context.Context, user *model.User) (uuid.UUID, error) {
	hashedPassword, err := s.hash.Hash(user.Password)
	if err != nil {
		s.log.Error("failed to hash password:", err)
		return uuid.Nil, err
	}

	u := model.User{
		UUID:     uuid.New(),
		Email:    user.Email,
		Password: hashedPassword,
	}

	userID, err := s.repo.Auth.Create(ctx, &u)
	if err != nil {
		s.log.Error("failed to create user:", err)
		return uuid.Nil, err
	}
	session := model.Session{
		ID: uuid.New(),
		RefreshToken: model.RefreshToken{
			UserID: userID,
		},
	}

	err = s.repo.Session.Create(ctx, session)
	if err != nil {
		s.log.Error("failed to create session:", err)
		return uuid.Nil, err
	}
	s.log.Info("user created successfully")

	return userID, nil
}

func (s *AuthService) Login(ctx context.Context, userID uuid.UUID, email, password string) (*model.AccessToken, *model.RefreshToken, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Error("failed to get user by credentials:", err)
		return nil, nil, err
	}

	if !s.hash.CompareHash(password, user.Password) {
		s.log.Error("invalid credentials")
		return nil, nil, errors.New("invalid password")
	}

	if user.UUID != userID {
		s.log.Error("user not found")
		return nil, nil, errors.New("user not found")
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

	session, err := s.repo.Session.GetByUserID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get session:", err)
		return nil, nil, err
	}
	ss := model.Session{
		ID: session.ID,
		RefreshToken: model.RefreshToken{
			ID:            refresh.ID,
			UserID:        userID,
			AccessTokenID: access.ID,
			Token:         hashedRefresh,
			ExpiresAt:     time.Now().Add(s.refreshTokenTTL),
		},
	}

	err = s.repo.Session.Update(ctx, ss)
	if err != nil {
		s.log.Error("failed to set session:", err)
		return nil, nil, err
	}

	s.log.Info("user logged in successfully")
	return access, refresh, nil
}

func (s *AuthService) Refresh(ctx context.Context, userID uuid.UUID, token string) (*model.AccessToken, *model.RefreshToken, error) {
	session, err := s.repo.Session.GetByUserID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get session:", err)
		return nil, nil, err
	}

	// Проверяем соответствие хешированного токена обновления переданному токену
	if !s.hash.CompareHash(token, session.RefreshToken.Token) {
		s.log.Error("failed to compare hash")
		return nil, nil, errors.New("failed to compare hash")
	}

	// Генерируем новый токен доступа
	access, refresh, err := s.jwt.GenerateTokenPair(userID, s.accessTokenTTL, s.refreshTokenTTL)
	if err != nil {
		s.log.Error("failed to generate token pair:", err)
		return nil, nil, err
	}

	// Хешируем новый токен обновления для сохранения в базе данных
	hashedRefresh, err := s.hash.Hash(refresh.Token)
	if err != nil {
		s.log.Error("failed to hash new refresh token:", err)
		return nil, nil, err
	}

	newSession := model.Session{
		ID: session.ID,
		RefreshToken: model.RefreshToken{
			ID:            refresh.ID,
			UserID:        userID,
			AccessTokenID: access.ID,
			Token:         hashedRefresh,
			ExpiresAt:     time.Now().Add(s.refreshTokenTTL),
		},
	}

	err = s.repo.Session.Update(ctx, newSession)
	if err != nil {
		s.log.Error("failed to update session:", err)
		return nil, nil, err
	}

	s.log.Info("token refreshed successfully")
	return access, refresh, nil
}
