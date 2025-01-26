package main

import (
	"context"
	"github.com/goocarry/rest-ultimate/internal/config"
	"github.com/goocarry/rest-ultimate/internal/http-server/handlers/user"
	mwLogger "github.com/goocarry/rest-ultimate/internal/http-server/middleware/logger"
	"github.com/goocarry/rest-ultimate/internal/lib/logger/sl"
	"github.com/goocarry/rest-ultimate/internal/storage/sqlite"
	"github.com/goocarry/rest-ultimate/internal/telegrambot"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting...", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("cannot create storage", sl.Err(err))
		os.Exit(1)
	}

	telegramBot, err := telegrambot.NewTelegramBot(cfg.TelegramBotToken, storage, log)
	if err != nil {
		log.Error("cannot create Telegram bot", sl.Err(err))
		os.Exit(1)
	}

	go func() {
		telegramBot.Start()
	}()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)

	router.Route("/user", func(r chi.Router) {

		r.Post("/register", user.New(log, storage))
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Info("received signal, shutting down...", slog.String("signal", sig.String()))

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info("shutdown complete")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	return log
}
