package protoclient

import (
	"fmt"

	pb "github.com/WhiteRain7/GoTestCase/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func Connect(cfg Config) (pb.ExchangeServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	client := pb.NewExchangeServiceClient(conn)
	return client, conn, nil
}

func MustConnect(cfg Config) (pb.ExchangeServiceClient, *grpc.ClientConn) {
	client, conn, err := Connect(cfg)
	if err != nil {
		panic(err)
	}
	return client, conn
}
