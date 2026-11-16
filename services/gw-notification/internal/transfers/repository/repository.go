package repository

import (
	"context"

	"github.com/WhiteRain7/GoTestCase/notification/internal/transfers/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{client.Database("main").Collection("transfers")}
}

func (r *Repository) Add(ctx context.Context, transfer domain.Transfer) error {
	doc := bson.M{
		"user":      transfer.User,
		"operation": transfer.Operation,
		"amount":    transfer.Amount,
		"currency":  transfer.Currency,
		"timestamp": transfer.Timestamp,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	return err
}
