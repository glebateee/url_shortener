package main

import (
	"log/slog"
	"os"
	"strings"
	"urlshortener/internal/config"
	mwLogger "urlshortener/internal/http-server/middleware/logger"
	"urlshortener/internal/lib/logger/sl"
	"urlshortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

/*
sqlite 		: github.com/mattn/go-sqlite3
config 		: github.com/ilyakaznacheev/cleanenv
chi-router  : github.com/go-chi/chi/v5
json parser : github.com/go-chi/render
validator 	: github.com/go-playground/validator/v10
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

	router := chi.NewRouter()

	// middleware

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
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
