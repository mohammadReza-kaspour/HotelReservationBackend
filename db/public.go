package db

import (
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db/shared"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	Client         *mongo.Client
	Database       *mongo.Database
	UserStorage    UserStorage
	HotelStorage   HotelStorage
	RoomStorage    RoomStorage
	RoleStorage    RoleStorage
	SessionStorage SessionStorage
	BookingStorage BookingStorage
}

func NewStorage() *Storage {
	client := shared.NewMongoDbClient()
	database := client.Database(config.DbName)

	return &Storage{
		Client:         client,
		Database:       database,
		UserStorage:    NewMongoUserStorage(database),
		HotelStorage:   NewMongoHotelStorage(database),
		RoomStorage:    NewMongoRoomStorage(database),
		RoleStorage:    NewMongoRoleStorage(database),
		SessionStorage: NewMongoSessionStorage(database),
		BookingStorage: NewMongoBookingStorage(database),
	}
}
