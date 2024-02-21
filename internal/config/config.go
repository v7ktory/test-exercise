package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultAccessTokenTTL  = 15 * time.Minute
	defaultRefreshTokenTTL = 24 * time.Hour * 30
)

type Config struct {
	Mongo MongoDB
	Auth  AuthConfig
}

type (
	MongoDB struct {
		URI      string
		Username string
		Password string
		DBName   string
	}
	AuthConfig struct {
		JWT          JWTConfig
		PasswordSalt string
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration
		RefreshTokenTTL time.Duration
		SigningKey      string
	}
)

func InitCfg() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var cfg Config
	setEnv(&cfg)

	cfg.Auth.JWT.AccessTokenTTL = defaultAccessTokenTTL
	cfg.Auth.JWT.RefreshTokenTTL = defaultRefreshTokenTTL

	return &cfg, nil
}

func setEnv(cfg *Config) {
	cfg.Mongo.URI = os.Getenv("MONGO_URI")
	cfg.Mongo.Username = os.Getenv("MONGO_USERNAME")
	cfg.Mongo.Password = os.Getenv("MONGO_PASS")

	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")
}
