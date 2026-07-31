package main

import (
	"os"
	"log/slog"
	"time"
	"context"

	"github.com/labstack/echo/v5"

	"github.com/nekto-sns/nekto-server/app/handler"
	"github.com/nekto-sns/nekto-server/app/repository"
	"github.com/nekto-sns/nekto-server/app/service"

	"github.com/nekto-sns/nekto-server/app/shared/config"
	"github.com/nekto-sns/nekto-server/app/shared/database"
	"github.com/nekto-sns/nekto-server/app/shared/errorhandler"
)

func main() {
	cfg := config.Load()

	logOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: false,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, logOpts))
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

	userRepo := repository.NewUserRepository(dbPool)
	userSvc  := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	e := echo.New()
	e.HTTPErrorHandler = errorhandler.ErrorHandler
	e.GET("/users/:name", userHandler.ByName)

	e.Start(cfg.Port)
}
