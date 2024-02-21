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

type JWTManager interface {
	GenerateTokenPair(userID uuid.UUID, accessTokenTTL, refreshTokenTTL time.Duration) (*model.AccessToken, *model.RefreshToken, error)
	generateAccessToken(userID uuid.UUID, ttl time.Duration) (*model.AccessToken, error)
	generateRefreshToken(userID uuid.UUID, ttl time.Duration) (*model.RefreshToken, error)
	validateToken(signedToken string) (string, error)
}

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
	refreshToken, err := j.generateRefreshToken(userID, refreshTokenTTL)
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

	accessToken := model.AccessToken{
		Token:  signedString,
		ID:     uuid.New(),
		UserID: userID,
	}
	return &accessToken, nil
}

func (j *JWT) generateRefreshToken(userID uuid.UUID, ttl time.Duration) (*model.RefreshToken, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return nil, err
	}
	token := base64.StdEncoding.EncodeToString(tokenBytes)

	refreshToken := model.RefreshToken{
		Token:  token,
		ID:     uuid.New(),
		UserID: userID,
	}
	return &refreshToken, nil
}

func (j *JWT) validateToken(signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("token validation failed")
	}

	return (*claims)["sub"].(string), nil
}
