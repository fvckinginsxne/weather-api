package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "weather-api/docs"
	"weather-api/internal/client/weathermap"
	"weather-api/internal/config"
	"weather-api/internal/lib/logger/sl"
	"weather-api/internal/service/weather"
	"weather-api/internal/storage/sqlite"
	"weather-api/internal/transport/handler/weather/create"
)

const (
	shutdownTimeout = 10 * time.Second
)

// @title Weather API
// @version 1.0
// @description API for getting current weather
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	cfg := config.MustLoad()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer cancel()

	weatherClient := weathermap.New(log, cfg.WeatherAPI.URL, cfg.WeatherAPI.Key)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}

	if err := storage.Init(); err != nil {
		panic(err)
	}

	weatherService := weather.New(log, weatherClient, storage)

	g := gin.New()

	g.Use(gin.Recovery())

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	weatherGroup := g.Group("/weather")
	{
		weatherGroup.POST("/", create.New(ctx, log, weatherService))
	}

	srvAddr := serverAddress(cfg)

	srv := &http.Server{
		Addr:         srvAddr,
		Handler:      g,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Debug("starting server", slog.String("address", srvAddr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("shutdown signal received")
	case err := <-serverErr:
		log.Error("server error", sl.Err(err))
		cancel()
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("failed to shutdown server", sl.Err(err))
	}
	if err := storage.Close(shutdownCtx); err != nil {
		log.Error("failed to close storage", sl.Err(err))
	}

	log.Info("shutdown complete")
}

func serverAddress(cfg *config.Config) string {
	return fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
}
