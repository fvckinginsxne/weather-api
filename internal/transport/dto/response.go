package dto

import "weather-api/internal/domain/model"

type ErrorResponse struct {
	Error string `json:"error" example:"invalid request"`
}

type WeatherResponse struct {
	City      string  `json:"city" example:"Нижний Новгород"`
	Desc      string  `json:"desc" example:"облачно"`
	Temp      float64 `json:"temp" example:"27.8"`
	WindSpeed float64 `json:"wind_speed" example:"1.79"`
}

func WeatherResponseToModel(weatherResponse *WeatherResponse) *model.Weather {
	return &model.Weather{
		City:      weatherResponse.City,
		Desc:      weatherResponse.Desc,
		Temp:      weatherResponse.Temp,
		WindSpeed: weatherResponse.WindSpeed,
	}
}

func WeatherModelToResponse(weather *model.Weather) *WeatherResponse {
	return &WeatherResponse{
		City:      weather.City,
		Desc:      weather.Desc,
		Temp:      weather.Temp,
		WindSpeed: weather.WindSpeed,
	}
}
