package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lekss361/curserv2/gateway/internal/dto"
	pb "github.com/lekss361/curserv2/pkg/currency"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CurrencyService interface {
	GetRatesByDate(ctx context.Context, date time.Time) (map[string]float64, error)
	GetRatesHistory(ctx context.Context, start, end time.Time) ([]dto.RateResponse, error)
}

type currencyService struct {
	client pb.CurrencyServiceClient
}

func NewCurrencyService(grpcURL string) (CurrencyService, error) {
	conn, err := grpc.Dial(grpcURL, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to currency service: %w", err)
	}
	client := pb.NewCurrencyServiceClient(conn)
	return &currencyService{client: client}, nil
}

func (s *currencyService) GetRatesByDate(ctx context.Context, date time.Time) (map[string]float64, error) {
	// Build and send the gRPC request
	req := &pb.GetRateByDateRequest{
		Date: timestamppb.New(date),
	}
	res, err := s.client.GetRateByDate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("rpc GetRateByDate failed: %w", err)
	}
	return res.GetRub(), nil
}

func (s *currencyService) GetRatesHistory(ctx context.Context, start, end time.Time) ([]dto.RateResponse, error) {
	req := &pb.GetHistoryRequest{
		Start: timestamppb.New(start),
		End:   timestamppb.New(end),
	}
	res, err := s.client.GetHistory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("rpc GetHistory failed: %w", err)
	}

	history := make([]dto.RateResponse, len(res.GetHistory()))
	for i, dr := range res.GetHistory() {
		history[i] = dto.RateResponse{
			Date: dr.GetDate().AsTime(),
			Rub:  dr.GetRub(),
		}
	}
	return history, nil
}
