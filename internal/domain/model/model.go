package model

import "time"

type Weather struct {
	City      string
	CreatedAt time.Time
	Desc      string
	Temp      float64
	WindSpeed float64
}
