package protoserver

import (
	"fmt"
	"net"

	grpcimp "github.com/WhiteRain7/GoTestCase/exchanger/internal/adapters/grpc"
	pb "github.com/WhiteRain7/GoTestCase/proto"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type Config struct {
	Port int `mapstructure:"port"`
}

func NewServer(exchangeService *grpcimp.ExchangeService) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)

	pb.RegisterExchangeServiceServer(grpcServer, exchangeService)

	grpc_prometheus.Register(grpcServer)

	return grpcServer
}

func Run(cfg Config, server *grpc.Server) {
	lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	server.Serve(lis)
}
