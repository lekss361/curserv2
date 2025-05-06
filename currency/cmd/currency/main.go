package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	configpkg "github.com/lekss361/curserv2/currency/internal/config"
	handlerpkg "github.com/lekss361/curserv2/currency/internal/handler"
	repo "github.com/lekss361/curserv2/currency/internal/repository"
	service "github.com/lekss361/curserv2/currency/services"
	worker "github.com/lekss361/curserv2/currency/worker"
	proto "github.com/lekss361/curserv2/pkg/currency"
	"go.uber.org/zap"
	"google.golang.org/grpc"

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

	// Open database connection
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

	// Create gRPC server and register service
	grpcServer := grpc.NewServer()
	h := handlerpkg.NewCurrencyHandler(ratesRepo)
	grpcHandler := handlerpkg.NewGRPCServer(h)
	proto.RegisterCurrencyServiceServer(grpcServer, grpcHandler)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	cw.Start(ctx)

	addr := fmt.Sprintf(":%d", cfg.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Failed to listen on gRPC port", zap.String("addr", addr), zap.Error(err))
	}
	logger.Info("gRPC server listening", zap.String("addr", addr))
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("gRPC serve error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("Shutdown signal received, stopping services...")

	grpcServer.GracefulStop()
	cw.Wait()
	logger.Info("Service stopped gracefully")
}
