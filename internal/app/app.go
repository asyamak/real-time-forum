package app

import (
	"os"
	"os/signal"
	"syscall"

	"real-time-forum/internal/config"
	handler "real-time-forum/internal/handler/http"
	"real-time-forum/internal/repository"
	"real-time-forum/internal/server"
	"real-time-forum/internal/service"
	"real-time-forum/pkg/hasher"
	"real-time-forum/pkg/logger"
	"real-time-forum/pkg/sqlite"
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

	db, err := sqlite.ConnectDatabase(cfg)
	if err != nil {
		a.log.Error("error while connecting database: %s", err.Error())
	}

	a.log.Info("Database connected")

	h, err := hasher.NewHasher("aboba")
	if err != nil {
		a.log.Error("error while connecting database: %s", err.Error())
	}

	repository := repository.NewRepository(db)
	service := service.NewService(repository, h)
	handler := handler.NewHandler(service)

	server := server.NewServer(cfg, handler.InitRoutes())

	quit := make(chan os.Signal, 1)

	go func() {
		a.log.Info("Starting server at port %v -> http://localhost%v", cfg.API.Port, cfg.API.Port)
		server.Start()
	}()

	a.log.Info("Real-Time-Forum app started")

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	select {
	case signal := <-quit:
		a.log.Info("signal accepted: %v", signal)
	case err := <-server.ServerErrNotify():
		a.log.Info("server closing: %v", err)
	}

	a.log.Info("Real-Time-Forum app shutting down...")

	if err := db.Close(); err != nil {
		a.log.Error(err.Error())
	}

	a.log.Info("Database closed")

	if err := server.Shutdown(); err != nil {
		a.log.Error(err.Error())
	}

	a.log.Info("Server stopped")
}
