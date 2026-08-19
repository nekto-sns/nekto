package main

import (
	"os"
	"log/slog"
	"net/http"
	"time"
	"context"

	"github.com/labstack/echo/v5"

	"github.com/redis/go-redis/v9"

	"github.com/nekto-sns/nekto/server/app/route"
	"github.com/nekto-sns/nekto/server/app/handler"
	"github.com/nekto-sns/nekto/server/app/repository"
	"github.com/nekto-sns/nekto/server/app/service"


	"github.com/nekto-sns/nekto/server/app/manager/sessionmanager"
	"github.com/nekto-sns/nekto/server/app/client/scratchauth"

	"github.com/nekto-sns/nekto/server/app/shared/config"
	"github.com/nekto-sns/nekto/server/app/shared/database"
	"github.com/nekto-sns/nekto/server/app/shared/errorhandler"
)

func main() {
	cfg := config.Load()

	logOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: false,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, logOpts))
	slog.SetDefault(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	dbPool, err := database.NewPool(ctx, cfg.DBUrl)
	defer dbPool.Close()
	if err != nil {
		slog.Error("Server startup failed", "error", err)
		os.Exit(1)
	}

	err = database.RunMigrations(ctx, dbPool)
	if err != nil {
		slog.Error("DB migration failed", "error", err)
		os.Exit(1)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 1000,
	})

	userRepo    := repository.NewUser(dbPool)
	userSvc     := service.NewUser(userRepo)
	userHandler := handler.NewUser(userSvc)

	sessionRepo := repository.NewSession(rdb)
	sessionMN   := sessionmanager.New(sessionRepo, 24 * 7 * time.Hour)

	saClient    := scratchauth.New(&http.Client{}, cfg.ScratchAuthURL, []string{cfg.LoginCallbackURL})
	authRepo    := repository.NewScratchAuth(dbPool)
	authSvc     := service.NewScratchAuth(authRepo, saClient, sessionMN)
	authHandler := handler.NewScratchAuth(authSvc, cfg.ScratchAuthURL, cfg.LoginCallbackURL)

	e := echo.New()
	e.HTTPErrorHandler = errorhandler.ErrorHandler

	route.Setup(e, userHandler, authHandler)


	e.Start(cfg.Port)
}
