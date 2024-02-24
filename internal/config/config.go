package config

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultAccessTokenTTL  = 30 * time.Minute
	defaultRefreshTokenTTL = 24 * time.Hour * 30 // 30 days

	defaultQueryTimeout = 10 * time.Second

	defaultPort           = "8080"
	defaultMaxHeaderBytes = 1 << 20
	defaultReadTimeout    = 10 * time.Second
	defaultWriteTimeout   = 10 * time.Second
)

type (
	Cfg struct {
		Mongo  MongoCfg
		Auth   AuthCfg
		Server Server
	}
	MongoCfg struct {
		Hosts        []string
		Username     string
		Password     string
		DB           string
		QueryTimeout time.Duration
	}
	AuthCfg struct {
		JWT          JWTCfg
		PasswordSalt string
	}
	JWTCfg struct {
		AccessTokenTTL  time.Duration
		RefreshTokenTTL time.Duration
		SigningKey      string
	}
	Server struct {
		Port           string
		MaxHeaderBytes int
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
	}
)

func InitCfg() (*Cfg, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	var cfg Cfg
	err = loadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	err = loadDefault(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadEnv(cfg *Cfg) error {
	mongoHosts := os.Getenv("MONGO_HOSTS")
	if mongoHosts == "" {
		return errors.New("missing MONGO_HOSTS")
	}

	cfg.Mongo.Hosts = strings.Split(mongoHosts, ",")
	cfg.Mongo.Username = os.Getenv("MONGO_USERNAME")
	cfg.Mongo.Password = os.Getenv("MONGO_PASS")
	cfg.Mongo.DB = os.Getenv("MONGO_DBNAME")

	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")

	return nil
}

func loadDefault(cfg *Cfg) error {
	cfg.Auth.JWT.AccessTokenTTL = defaultAccessTokenTTL
	cfg.Auth.JWT.RefreshTokenTTL = defaultRefreshTokenTTL

	cfg.Mongo.QueryTimeout = defaultQueryTimeout

	cfg.Server.Port = defaultPort
	cfg.Server.MaxHeaderBytes = defaultMaxHeaderBytes
	cfg.Server.ReadTimeout = defaultReadTimeout
	cfg.Server.WriteTimeout = defaultWriteTimeout

	return nil
}
