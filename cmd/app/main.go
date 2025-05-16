package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"weather-api/internal/config"
	"weather-api/internal/lib/logger/sl"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	log.Debug("config loaded", slog.Any("config", cfg))

	g := gin.New()

	srv := &http.Server{
		Addr:         serverAddress(cfg),
		Handler:      g,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	log.Info("starting server", slog.String("addr", srv.Addr))

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server error", sl.Err(err))
	}

	log.Info("server shutdown")
}

func serverAddress(cfg *config.Config) string {
	return fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
}
