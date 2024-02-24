package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/v7ktory/test/internal/model"
)

type JWT struct {
	signingKey string
}

func NewJWT(signingKey string) *JWT {
	return &JWT{
		signingKey: signingKey,
	}
}

func (j *JWT) GenerateTokenPair(userID uuid.UUID, accessTokenTTL, refreshTokenTTL time.Duration) (*model.AccessToken, *model.RefreshToken, error) {
	accessToken, err := j.generateAccessToken(userID, accessTokenTTL)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := j.generateRefreshToken(userID, accessToken.ID, refreshTokenTTL)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (j *JWT) generateAccessToken(userID uuid.UUID, ttl time.Duration) (*model.AccessToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
	})

	signedString, err := token.SignedString([]byte(j.signingKey))
	if err != nil {
		return nil, err
	}

	accessToken := &model.AccessToken{
		ID:     uuid.New(),
		UserID: userID,
		Token:  signedString,
	}
	return accessToken, nil
}

func (j *JWT) generateRefreshToken(userID, accessTokenID uuid.UUID, ttl time.Duration) (*model.RefreshToken, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return nil, err
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	refreshToken := &model.RefreshToken{
		ID:            uuid.New(),
		UserID:        userID,
		AccessTokenID: accessTokenID,
		Token:         token,
		ExpiresAt:     time.Now().Add(ttl),
	}
	return refreshToken, nil
}

func (j *JWT) ValidateToken(signedToken string) (*uuid.UUID, error) {
	// Парсинг и валидация токена
	token, err := jwt.ParseWithClaims(signedToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.signingKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	// Проверка наличия и валидности токена
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token validation failed")
	}

	// Получение идентификатора пользователя из токена
	subject, ok := (*claims)["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("subject not found in claims")
	}

	// Парсинг идентификатора пользователя в UUID
	userID, err := uuid.Parse(subject)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return &userID, nil
}
