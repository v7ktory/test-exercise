package app

import (
	"log"

	"github.com/v7ktory/test/internal/config"
	"github.com/v7ktory/test/pkg/database/mongodb"
)

func Run() {
	cfg, err := config.InitCfg()
	if err != nil {
		log.Fatal(err)
	}
	mongo, err := mongodb.NewMongoDB(cfg.Mongo.URI, cfg.Mongo.Username, cfg.Mongo.Password)
	if err != nil {
		log.Fatal(err)
	}
}
