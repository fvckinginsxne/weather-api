package weathermap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-resty/resty/v2"

	"weather-api/internal/client"
	"weather-api/internal/transport/dto"
)

var (
	ErrCityNotFound = errors.New("city not found")
)

type Client struct {
	log    *slog.Logger
	client *resty.Client
	url    string
	key    string
}

func New(
	log *slog.Logger,
	url, key string,
) *Client {
	return &Client{
		log:    log,
		client: resty.New(),
		url:    url,
		key:    key,
	}
}

type Response struct {
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
		Pressure  int     `json:"pressure"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

func (c *Client) Weather(ctx context.Context, city string) (*dto.WeatherResponse, error) {
	const op = "client.weathermap.Weather"

	log := c.log.With(
		slog.String("op", op),
		slog.String("city", city),
	)

	log.Info("fetching weather")

	ctx, cancel := context.WithTimeout(context.Background(), client.APIRequestTimeout)
	defer cancel()

	var result *Response
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&result).
		SetQueryParam("q", city).
		SetQueryParam("appid", c.key).
		Get(c.url)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("%s: %w", op, ctx.Err())
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("openweathermap response status", slog.String("status", resp.Status()))

	if resp.StatusCode() == http.StatusNotFound {
		return nil, fmt.Errorf("%s: %w", op, ErrCityNotFound)
	}

	log.Info("weather successfully fetched from OpenWeatherMap")

	return &dto.WeatherResponse{
		City:      city,
		Desc:      result.Weather[0].Description,
		Temp:      result.Main.Temp,
		WindSpeed: result.Wind.Speed,
	}, nil
}
