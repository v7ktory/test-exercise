package app

import (
	"log"

	"github.com/v7ktory/test/internal/config"
	"github.com/v7ktory/test/internal/repository"
	"github.com/v7ktory/test/internal/service"
	"github.com/v7ktory/test/pkg/database/mongodb"
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
	repository := repository.NewRepository(db)
	service := service.NewService(
		repository,
		cfg.Auth.PasswordSalt,
		cfg.Auth.JWT.SigningKey,
		log,
		cfg.Auth.JWT.AccessTokenTTL,
	)
}
