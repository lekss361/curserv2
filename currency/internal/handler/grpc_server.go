package handler

import (
	"context"

	pb "github.com/lekss361/curserv2/pkg/currency"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	h *CurrencyHandlerImpl
	pb.UnimplementedCurrencyServiceServer
}

func NewGRPCServer(h *CurrencyHandlerImpl) *Server {
	return &Server{h: h}
}

func (s *Server) GetRateByDate(
	ctx context.Context,
	req *pb.GetRateByDateRequest,
) (*pb.GetRateByDateResponse, error) {
	t := req.GetDate().AsTime()
	rates, err := s.h.GetRateByDate(ctx, t)
	if err != nil {
		return nil, err
	}
	return &pb.GetRateByDateResponse{
		Date: timestamppb.New(t),
		Rub:  rates,
	}, nil
}

// GetHistory handles the GetHistory gRPC call
func (s *Server) GetHistory(
	ctx context.Context,
	req *pb.GetHistoryRequest,
) (*pb.GetHistoryResponse, error) {
	start := req.GetStart().AsTime()
	end := req.GetEnd().AsTime()
	history, err := s.h.GetHistory(ctx, start, end)
	if err != nil {
		return nil, err
	}
	resp := make([]*pb.DateRates, len(history))
	for i, dr := range history {
		resp[i] = &pb.DateRates{
			Date: timestamppb.New(dr.Date),
			Rub:  dr.Rates,
		}
	}
	return &pb.GetHistoryResponse{History: resp}, nil
}
