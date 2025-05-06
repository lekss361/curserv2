package dto

import "time"

type RateResponse struct {
	Date time.Time          `json:"date"`
	Rub  map[string]float64 `json:"rub"`
}

type HistoryResponse struct {
	History []RateResponse `json:"history"`
}
