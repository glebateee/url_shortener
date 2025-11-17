package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"urlshortener/internal/config"
	"urlshortener/internal/http-server/handlers/url/redirect"
	"urlshortener/internal/http-server/handlers/url/save"
	mwLogger "urlshortener/internal/http-server/middleware/logger"
	"urlshortener/internal/lib/logger/env"
	"urlshortener/internal/lib/logger/sl"
	"urlshortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	// init config
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// init logger
	log := setupLogger(cfg.Env)
	log = log.With(
		slog.String("env", cfg.Env),
	)
	log.Debug("debug")
	log.Info("info")

	//init storage

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	log.Debug("storage init successfully")

	id, err := storage.SaveURL("a", "b")
	if err != nil {
		log.Debug("failed to init storage", sl.Err(err))
	}
	log.Debug("added new alias", slog.Int64("id", id))
	id, err = storage.SaveURL("a", "b")
	if err != nil {
		log.Debug("failed to init storage", sl.Err(err))
	}
	log.Debug("added new alias", slog.Int64("id", id))

	// init router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url_shortener", map[string]string{
			cfg.HttpServer.User: cfg.HttpServer.Password,
		}))
		r.Post("/", save.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))
	log.Info("starting server", slog.String("address", cfg.Adress))

	server := &http.Server{
		Addr:         cfg.Adress,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	if err = server.ListenAndServe(); err != nil {
		log.Error("failed starting server", sl.Err(err))
	}

	log.Error("server stopped")
}

func setupLogger(environment string) *slog.Logger {
	var log *slog.Logger
	switch strings.ToLower(environment) {
	case env.Local:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case env.Dev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case env.Prod:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
