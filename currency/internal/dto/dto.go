package dto

import (
	"time"
)

type RatesResponse struct {
	Date time.Time          `json:"date"` //  "2025-04-26"
	Rub  map[string]float64 `json:"rub"`  // код валюты курс
}

type RatesRepo interface {
	Save(date time.Time, rates map[string]float64) error
	Get(date time.Time) (map[string]float64, error)
}
