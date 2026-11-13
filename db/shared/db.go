package shared

import (
	"context"
	"log"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Functions
func NewMongoDbClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI(config.DbUri).SetTimeout(config.DbConnectionTimout)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func SizeAndPageInitializer(size **int64, page *	*int64) {
	initSize := int64(config.DefaultSize)
	initPage := int64(config.DefaultPage)

	if *page == nil {
		*page = &initPage
	}
	if *size == nil {
		*size = &initSize
	}
}


// Types
type Dropper interface {
	Drop(ctx context.Context) error
}
type Statistician interface {
	Info(ctx context.Context) (*StorageInfoSchema, error)
}

type StorageInfoSchema struct {
	Counts int
}
