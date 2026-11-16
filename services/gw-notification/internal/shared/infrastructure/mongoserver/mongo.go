package mongoserver

import (
	"context"
	"fmt"

	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"password"`
	Base string `mapstructure:"base"`
}

func Connect(cfg Config) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Base))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func MustConnect(cfg Config) *mongo.Client {
	mongodb, err := Connect(cfg)
	if err != nil {
		panic(err)
	}
	return mongodb
}
