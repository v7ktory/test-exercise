package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager interface {
	GenerateTokenPair(userID string, ttl time.Duration) (string, string, error)
	generateAccessToken(userID string, ttl time.Duration) (string, error)
	generateRefreshToken() (string, error)
	validateToken(signedToken string) (string, error)
}

type JWT struct {
	secretKey string
}

func NewJWT(secretKey string) *JWT {
	return &JWT{
		secretKey: secretKey,
	}
}

func (j *JWT) GenerateTokenPair(userID string, ttl time.Duration) (string, string, error) {

	access, err := j.generateAccessToken(userID, ttl)
	if err != nil {
		return "", "", err
	}

	refresh, err := j.generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (j *JWT) generateAccessToken(userID string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
	})

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWT) generateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	token := base64.StdEncoding.EncodeToString(tokenBytes)

	return token, nil
}

func (j *JWT) validateToken(signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
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
