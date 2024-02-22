package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultAccessTokenTTL  = 30 * time.Minute
	defaultRefreshTokenTTL = 24 * time.Hour * 30

	defaultQueryTimeout = 5 * time.Second

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
	loadEnv(&cfg)

	loadDefault(&cfg)
	return &cfg, nil
}

func loadEnv(cfg *Cfg) {
	mongoHosts := os.Getenv("MONGO_HOSTS")

	cfg.Mongo.Hosts = strings.Split(mongoHosts, ",")
	cfg.Mongo.Username = os.Getenv("MONGO_USERNAME")
	cfg.Mongo.Password = os.Getenv("MONGO_PASS")

	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")

}
func loadDefault(cfg *Cfg) {
	cfg.Auth.JWT.AccessTokenTTL = defaultAccessTokenTTL
	cfg.Auth.JWT.RefreshTokenTTL = defaultRefreshTokenTTL

	cfg.Mongo.QueryTimeout = defaultQueryTimeout

	cfg.Server.Port = defaultPort
	cfg.Server.MaxHeaderBytes = defaultMaxHeaderBytes
	cfg.Server.ReadTimeout = defaultReadTimeout
	cfg.Server.WriteTimeout = defaultWriteTimeout
}
