package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/v7ktory/test/internal/config"
	"github.com/v7ktory/test/internal/repository"
	"github.com/v7ktory/test/internal/server"
	"github.com/v7ktory/test/internal/service"
	"github.com/v7ktory/test/internal/transport/rest"
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
	// todo конфиг допиши
	mongo, err := mongodb.NewMongoDB(context.TODO(), cfg.Mongo)
	if err != nil {
		log.Fatal(err)
	}

	log := logger.NewLogger()
	hash := hash.NewHasher(cfg.Auth.PasswordSalt)
	jwt := jwt.NewJWT(cfg.Auth.JWT.SigningKey)
	accessTTL := cfg.Auth.JWT.AccessTokenTTL
	refreshTTL := cfg.Auth.JWT.RefreshTokenTTL
	repository := repository.NewRepository(mongo)
	service := service.NewService(
		*repository,
		*hash,
		*jwt,
		log,
		accessTTL,
		refreshTTL,
	)
	handler := rest.NewHandler(*service)
	srv := server.NewServer(cfg, handler.InitRoutes())

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Error("error occurred while running http server: %s", err)

		}
	}()

	log.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Error("failed to stop server: %v", err)
	}

	if err := mongo.GetClient().Disconnect(context.Background()); err != nil {
		log.Error(err.Error())
	}
}
