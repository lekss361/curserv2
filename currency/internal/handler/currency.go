package handler

import (
	"context"
	"fmt"
	"github.com/lekss361/curserv2/currency/internal/dto"
	"time"
)

type CurrencyHandlerImpl struct {
	repo dto.RatesRepo
}

func NewCurrencyHandler(repo dto.RatesRepo) *CurrencyHandlerImpl {
	return &CurrencyHandlerImpl{repo: repo}
}

func (h *CurrencyHandlerImpl) GetRateByDate(ctx context.Context, date time.Time) (map[string]float64, error) {
	rates, err := h.repo.Get(date)
	if err != nil {
		return nil, fmt.Errorf("failed to get rates for date %s: %w", date.Format("2006-01-02"), err)
	}
	return rates, nil
}

func (h *CurrencyHandlerImpl) GetHistory(ctx context.Context, start, end time.Time) ([]DateRates, error) {
	if end.Before(start) {
		return nil, fmt.Errorf("end date must not be before start date")
	}

	var history []DateRates
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		rates, err := h.repo.Get(d)
		if err != nil {
			return nil, fmt.Errorf("failed to get rates for date %s: %w", d.Format("2006-01-02"), err)
		}
		history = append(history, DateRates{
			Date:  d,
			Rates: rates,
		})
	}

	return history, nil
}
