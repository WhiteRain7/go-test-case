package grpc

import (
	"context"

	"github.com/WhiteRain7/GoTestCase/exchanger/internal/rates/domain"
	pb "github.com/WhiteRain7/GoTestCase/proto"
)

type ExchangeService struct {
	pb.UnimplementedExchangeServiceServer
	uc UseCase
}

type UseCase interface {
	Get(context.Context, string) (*domain.Rate, error)
	All(context.Context) (map[string]float32, error)
}

func NewExchangeService(uc UseCase) *ExchangeService {
	return &ExchangeService{uc: uc}
}

func (s *ExchangeService) GetRates(ctx context.Context, req *pb.Empty) (*pb.GetRatesResponse, error) {
	rates, err := s.uc.All(ctx)
	if err != nil {
		return &pb.GetRatesResponse{}, err
	}

	return &pb.GetRatesResponse{Rates: rates}, nil
}

func (s *ExchangeService) Exchange(ctx context.Context, req *pb.ExchangeRequest) (*pb.ExchangeResponse, error) {
	rateFrom, err := s.uc.Get(ctx, req.FromCurrency)
	if err != nil {
		return &pb.ExchangeResponse{}, err
	}
	rateTo, err := s.uc.Get(ctx, req.ToCurrency)
	if err != nil {
		return &pb.ExchangeResponse{}, err
	}

	return &pb.ExchangeResponse{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         rateFrom.Rate / rateTo.Rate,
	}, nil
}
