package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	configpkg "github.com/lekss361/curserv2/currency/internal/config"
	repo "github.com/lekss361/curserv2/currency/internal/repository"
	service "github.com/lekss361/curserv2/currency/services"
	worker "github.com/lekss361/curserv2/currency/worker"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	cfg, err := configpkg.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Подключение к базе данных
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to open database", zap.Error(err))
	}
	defer db.Close()

	ratesRepo := repo.NewRatesRepo(db)
	svc := service.NewCurrencyService(ratesRepo, logger)

	interval := 24 * time.Hour
	baseCurrency := os.Getenv("BASE_CURRENCY")
	if baseCurrency == "" {
		baseCurrency = "RUB"
	}
	cw := worker.NewCurrencyWorker(svc, logger, interval, baseCurrency)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	cw.Start(ctx)

	<-ctx.Done()
	logger.Info("Shutdown signal received, stopping currency worker...")

	// Остановка воркера
	cw.Wait()
	logger.Info("Currency worker stopped")
}
