package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"weather-api/internal/domain/model"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Init() error {
	const op = "storage.sqlite.Init"

	q, err := s.db.Prepare(`
		CREATE TABLE IF NOT EXISTS weather_info (
		    id INTEGER PRIMARY KEY,
		    city VARCHAR(255) NOT NULL,
		    description VARCHAR(255) NOT NULL,
		    temperature FLOAT NOT NULL,
		    wind_speed FLOAT NOT NULL,
		    created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := q.Exec(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveWeather(ctx context.Context, weather *model.Weather) error {
	const op = "storage.sqlite.SaveWeather"

	q, err := s.db.PrepareContext(ctx, `
		INSERT INTO weather_info (city, description, temperature, wind_speed, created_at)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = q.ExecContext(
		ctx,
		weather.City,
		weather.Desc,
		weather.Temp,
		weather.WindSpeed,
		weather.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Close(cxt context.Context) error {
	done := make(chan struct{})

	var closeErr error
	go func() {
		closeErr = s.db.Close()
		close(done)
	}()

	select {
	case <-done:
		return closeErr
	case <-cxt.Done():
		return cxt.Err()
	}
}
