package worker

import (
	"context"
	service "github.com/lekss361/curserv2/currency/services"
	"go.uber.org/zap"
	"time"
)

// CurrencyWorker выполняет периодическое обновление курсов
type CurrencyWorker struct {
	service      service.CurrencyService
	logger       *zap.Logger
	interval     time.Duration
	baseCurrency string
	shutdownChan chan struct{}
}

// NewCurrencyWorker создает новый воркер
func NewCurrencyWorker(service service.CurrencyService, logger *zap.Logger, interval time.Duration, baseCurrency string) *CurrencyWorker {
	return &CurrencyWorker{
		service:      service,
		logger:       logger,
		interval:     interval,
		baseCurrency: baseCurrency,
		shutdownChan: make(chan struct{}),
	}
}

// Start запускает воркер в отдельной горутине
func (w *CurrencyWorker) Start(ctx context.Context) {
	w.logger.Info("Starting currency worker", zap.Duration("interval", w.interval))

	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				w.logger.Info("Currency worker received shutdown signal (ctx)")
				close(w.shutdownChan)
				return
			case <-ticker.C:
				w.logger.Info("Currency worker: fetching and saving rates...")
				if err := w.service.FetchAndSaveRates(ctx, w.baseCurrency); err != nil {
					w.logger.Error("Currency worker failed to fetch and save rates", zap.Error(err))
				} else {
					w.logger.Info("Currency worker: rates updated successfully")
				}
			}
		}
	}()
}

// Wait blocks until worker is shutdown
func (w *CurrencyWorker) Wait() {
	<-w.shutdownChan
}
