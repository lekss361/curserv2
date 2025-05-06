package handler

import (
	"context"
	"time"
)

type CurrencyHandler interface {
	GetRateByDate(ctx context.Context, date time.Time) (map[string]float64, error)
	GetHistory(ctx context.Context, start, end time.Time) ([]DateRates, error)
}

type DateRates struct {
	Date  time.Time
	Rates map[string]float64
}
