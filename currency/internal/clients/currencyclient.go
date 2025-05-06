package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lekss361/curserv2/currency/internal/config"
	"github.com/lekss361/curserv2/currency/internal/dto"
	"log"
	"net/http"
	"time"
)

func GetCurs(ctx context.Context) (*dto.RatesResponse, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	resp, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.ExternalServiceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("http.Get error: %w", err)
	}
	defer resp.Body.Close()

	var tmp struct {
		Date string             `json:"date"`
		Rub  map[string]float64 `json:"rub"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tmp); err != nil {
		return nil, fmt.Errorf("json.Decode error: %w", err)
	}

	const layout = "2006-01-02"
	parsedDate, err := time.Parse(layout, tmp.Date)
	if err != nil {
		return nil, fmt.Errorf("time.Parse error: %w", err)
	}

	result := &dto.RatesResponse{
		Date: parsedDate,
		Rub:  tmp.Rub,
	}
	return result, nil
}
