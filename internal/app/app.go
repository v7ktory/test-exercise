package app

import (
	"log"

	"github.com/v7ktory/test/internal/config"
	"github.com/v7ktory/test/internal/repository"
	"github.com/v7ktory/test/internal/service"
	"github.com/v7ktory/test/pkg/database/mongodb"
	"github.com/v7ktory/test/pkg/hash"
	"github.com/v7ktory/test/pkg/jwt"
	"github.com/v7ktory/test/pkg/logger"
)

func Run() {
	cfg, err := config.InitCfg()
	if err != nil {
		log.Fatal(err)
	}
	mongoClient, err := mongodb.NewMongoDB(cfg.Mongo.URI, cfg.Mongo.Username, cfg.Mongo.Password)
	if err != nil {
		log.Fatal(err)
	}
	db := mongoClient.Database(cfg.Mongo.DBName)
	log := logger.NewLogger()
	hash := hash.NewHasher(cfg.Auth.PasswordSalt)
	jwt := jwt.NewJWT(cfg.Auth.JWT.SigningKey)
	accessTTL := cfg.Auth.JWT.AccessTokenTTL
	refreshTTL := cfg.Auth.JWT.RefreshTokenTTL
	repository := repository.NewRepository(db)
	service := service.NewService(
		*repository,
		*hash,
		*jwt,
		log,
		accessTTL,
		refreshTTL,
	)

}
