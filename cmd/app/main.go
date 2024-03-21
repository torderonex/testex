package main

import (
	"log/slog"
	"os"
	"testex/internal/config"
	"testex/internal/handler"
	"testex/internal/service"
	"testex/internal/storage"
	"testex/internal/storage/postgres"
	"testex/pkg/server"
	"testex/pkg/slog"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()
	//logger init
	logger := setupLogger(cfg.Env)
	logger.Info("App is starting on port 8080", slog.String("Env", cfg.Env))

	//pg init
	db, err := postgres.New(cfg.Postgres)
	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	//app storage init
	appStorage := storage.New(db)

	//service init
	services := service.New(appStorage, logger, cfg)
	//router init
	router := handler.New(services, logger)
	_ = router
	//server init
	srv := server.New(cfg.HTTPServer.Port, router.Mux, cfg.HTTPServer.Timeout)
	err = srv.Run()
	if err != nil {
		logger.Error("failed to start server", sl.Err(err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
