package weather

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"weather-api/internal/client/weathermap"
	"weather-api/internal/domain/model"
	"weather-api/internal/lib/logger/sl"
	"weather-api/internal/transport/dto"
)

var (
	ErrCityNotFound = errors.New("city not found")
)

type Provider interface {
	Weather(ctx context.Context, city string) (*dto.WeatherResponse, error)
}

type Saver interface {
	SaveWeather(ctx context.Context, weather *model.Weather) error
}

type Service struct {
	log      *slog.Logger
	provider Provider
	saver    Saver
}

func New(
	log *slog.Logger,
	provider Provider,
	saver Saver,
) *Service {
	return &Service{
		log:      log,
		provider: provider,
		saver:    saver,
	}
}

func (s *Service) Save(
	ctx context.Context,
	createRequest *dto.CreateRequest,
) (*dto.WeatherResponse, error) {
	const op = "service.weather.Save"

	log := s.log.With(slog.String("op", op))

	weatherResponse, err := s.provider.Weather(ctx, createRequest.City)
	if err != nil {
		log.Error("failed to fetch weather", sl.Err(err))

		if errors.Is(err, weathermap.ErrCityNotFound) {

			return nil, fmt.Errorf("%s: %w", op, ErrCityNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	weather := dto.WeatherResponseToModel(weatherResponse)

	weather.CreatedAt = time.Now()

	log.Info("saving weather")

	if err = s.saver.SaveWeather(ctx, weather); err != nil {
		log.Error("failed to save weather", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("weather saved successfully")

	return dto.WeatherModelToResponse(weather), nil
}
