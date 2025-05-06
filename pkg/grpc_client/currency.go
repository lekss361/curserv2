package grpc_client

import (
	"context"
	"fmt"
	pb "github.com/lekss361/curserv2/pkg/currency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.CurrencyServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect gRPC: %w", err)
	}
	return &Client{
		conn:   conn,
		client: pb.NewCurrencyServiceClient(conn)}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetRateByDate(ctx context.Context, date time.Time) (map[string]float64, error) {
	req := &pb.GetRateByDateRequest{
		Date: timestamppb.New(date),
	}

	resp, err := c.client.GetRateByDate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC GetRateByDate failed: %w", err)
	}

	return resp.Rub, nil
}

func (c *Client) GetHistory(ctx context.Context, start, end time.Time) ([]*pb.DateRates, error) {
	req := &pb.GetHistoryRequest{
		Start: timestamppb.New(start),
		End:   timestamppb.New(end),
	}

	resp, err := c.client.GetHistory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC GetHistory failed: %w", err)
	}

	return resp.History, nil
}
