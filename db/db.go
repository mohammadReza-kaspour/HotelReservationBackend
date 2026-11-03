package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbUri           = "mongodb://localhost:27017"
	userCollection  = "users"
	hotelCollection = "hotels"
	roomCollection  = "rooms"
)

func NewMongoDbClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI(dbUri).SetTimeout(time.Second * 10)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

type Dropper interface {
	Drop(ctx context.Context) error
}
