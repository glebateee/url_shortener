package main

import (
	"log/slog"
	"os"
	"strings"
	"urlshortener/internal/config"
	"urlshortener/internal/lib/logger/sl"
	"urlshortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

/*
sqlite : github.com/mattn/go-sqlite3
config : github.com/ilyakaznacheev/cleanenv
*/
func main() {
	// Init config :
	cfg := config.MustLoad()
	// init logger : slog
	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))

	log.Debug("debug invoked")
	// init storage : sqlite

	_, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	// init router : chi

	// run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch strings.ToLower(env) {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
