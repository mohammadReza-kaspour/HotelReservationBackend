package shared

import (
	"context"
	"log"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDbClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI(config.DbUri).SetTimeout(config.DbConnectionTimout)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

type Dropper interface {
	Drop(ctx context.Context) error
}
