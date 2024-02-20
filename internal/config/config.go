package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Mongo MongoDB
}

type MongoDB struct {
	URI      string
	Username string
	Password string
}

func InitCfg() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		MongoDB{
			URI:      os.Getenv("MONGO_URI"),
			Username: os.Getenv("MONGO_USERNAME"),
			Password: os.Getenv("MONGO_PASSWORD"),
		},
	}, nil
}
