package create

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"weather-api/internal/lib/logger/sl"
	weatherService "weather-api/internal/service/weather"
	"weather-api/internal/transport/dto"
)

type WeatherSaver interface {
	Save(ctx context.Context, request *dto.CreateRequest) (*dto.WeatherResponse, error)
}

func New(
	ctx context.Context,
	log *slog.Logger,
	weatherSaver WeatherSaver,
) gin.HandlerFunc {
	const op = "handler.weather.create.New"

	return func(c *gin.Context) {
		log = log.With(slog.String("op", op))

		var req dto.CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")

				c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "request body is empty"})
				return
			}
			log.Error("failed to decode request body", sl.Err(err))

			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
			return
		}

		log.Debug("request body decoded", slog.Any("request", req))

		weather, err := weatherSaver.Save(ctx, &req)
		if err != nil {
			if errors.Is(err, weatherService.ErrCityNotFound) {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "city not found"})

				return
			}

			log.Error("failed to save weather", sl.Err(err))

			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
			return
		}

		c.JSON(http.StatusCreated, weather)
	}
}
