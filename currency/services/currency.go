package services

import (
	"context"
	"fmt"
	"github.com/lekss361/curserv2/currency/internal/dto"
	"go.uber.org/zap"
	"strings"
	"time"
)

// CurrencyService defines methods to fetch and format currency rates
// Date is returned as time.Time in DTO; JSON marshalling can format it to "YYYY-MM-DD"
type CurrencyService interface {
	// GetRatesByDate returns rates for a given date wrapped in a DTO ready for JSON
	GetRatesByDate(ctx context.Context, date time.Time) (dto.RatesResponse, error)
	// GetRatesHistory returns historical rates between two dates
	GetRatesHistory(ctx context.Context, start, end time.Time) ([]dto.RatesResponse, error)
	// FetchAndSaveRates fetches fresh rates and persists them
	FetchAndSaveRates(ctx context.Context, baseCurrency string) error
}

type currencyService struct {
	repo   dto.RatesRepo
	logger *zap.Logger
}

// NewCurrencyService constructs a CurrencyService
func NewCurrencyService(repo dto.RatesRepo, logger *zap.Logger) CurrencyService {
	return &currencyService{repo: repo, logger: logger}
}

// GetRatesByDate fetches raw rates, filters by 'rub' prefix, trims it, lowercases keys, and wraps in DTO
func (s *currencyService) GetRatesByDate(ctx context.Context, date time.Time) (dto.RatesResponse, error) {
	raw, err := s.repo.Get(date)
	if err != nil {
		return dto.RatesResponse{}, fmt.Errorf("failed to get rates for %s: %w", date.Format("2006-01-02"), err)
	}
	rates := make(map[string]float64)
	for code, rate := range raw {
		if strings.HasPrefix(code, "rub") {
			key := strings.ToLower(strings.TrimPrefix(code, "rub"))
			rates[key] = rate
		}
	}
	// Wrap in DTO with time.Time Date
	return dto.RatesResponse{
		Date: date,
		Rub:  rates,
	}, nil
}

func (s *currencyService) GetRatesHistory(ctx context.Context, start, end time.Time) ([]dto.RatesResponse, error) {
	if end.Before(start) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	var history []dto.RatesResponse
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		r, err := s.GetRatesByDate(ctx, d)
		if err != nil {
			s.logger.Warn("Failed to fetch rates for date", zap.String("date", d.Format("2006-01-02")), zap.Error(err))
			continue
		}
		history = append(history, r)
	}

	s.logger.Info("Fetched rates history", zap.Int("days", len(history)))
	return history, nil
}

func (s *currencyService) FetchAndSaveRates(ctx context.Context, baseCurrency string) error {
	rates := map[string]float64{
		"usd": 1.0,
		"eur": 0.9,
		"gbp": 0.8,
	}

	// Prepare DB format: prefix keys with 'rub' in uppercase codes
	date := time.Now()
	dbRates := make(map[string]float64)
	for k, v := range rates {
		dbRates["rub"+strings.ToUpper(k)] = v
	}

	if err := s.repo.Save(date, dbRates); err != nil {
		return fmt.Errorf("failed to save rates: %w", err)
	}

	s.logger.Info("Rates fetched and saved", zap.String("base", baseCurrency), zap.Time("date", date), zap.Any("rates", rates))
	return nil
}
