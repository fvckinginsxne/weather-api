package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server      ServerConfig     `env-prefix:"SERVER_" env-required:"true"`
	StoragePath string           `env:"STORAGE_PATH" env-required:"true"`
	WeatherAPI  WeatherAPIConfig `env-prefix:"WEATHER_API_" env-required:"true"`
}

type ServerConfig struct {
	Host        string        `env:"HOST" env-default:"localhost"`
	Port        int           `env:"PORT" env-default:"8080"`
	Timeout     time.Duration `env:"TIMEOUT" env-default:"30s"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
}

type WeatherAPIConfig struct {
	URL string `env:"URL" env-required:"true"`
	Key string `env:"KEY" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
