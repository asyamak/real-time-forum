package app

import (
	"context"
	"os"
	"os/signal"
	"real-time-forum/internal/config"
	handler "real-time-forum/internal/handler/http"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/server"
	"real-time-forum/internal/service"
	"real-time-forum/pkg/database"
	"real-time-forum/pkg/logger"
	"syscall"
)

type App struct {
	log *logger.Logger
}

func New() *App {
	return &App{
		log: logger.NewLogger("[Forum]"),
	}
}

func (a *App) Start(configPath *string, databaseName string) {
	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		a.log.Error(err.Error())
	}

	a.log.Info("Configs initialized")

	db, err := database.New(databaseName).ConnectDatabase(cfg)
	if err != nil {
		a.log.Error("error while connecting database: %s", err.Error())
	}

	a.log.Info("Database connected")

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)

	server := server.NewServer(cfg, handler.InitRoutes())

	a.log.Info("Starting server at port %v -> http://localhost%v", cfg.API.Port, cfg.API.Port)

	go func() {
		if err := server.Start(); err != nil {
			a.log.Info(err.Error())
		}
	}()

	a.log.Info("Real-Time-Forum app started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	a.log.Info("Real-Time-Forum app shutting down...")

	if err := db.Close(); err != nil {
		a.log.Error(err.Error())
	}

	a.log.Info("Database closed")

	if err := server.Shutdown(context.Background()); err != nil {
		a.log.Error(err.Error())
	}

	a.log.Info("Server stopped")
}
