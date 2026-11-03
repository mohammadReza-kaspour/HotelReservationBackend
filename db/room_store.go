package db

import (
	"context"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	Dropper

	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	/*
		todo!! => this is piece of shit and its anti pettern and
		must replce with service layer and exists only for educational perpose
	*/
	HotelStore

	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore, dbName string) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(dbName).Collection(roomCollection),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) Drop(ctx context.Context) error {
	return s.collection.Drop(ctx)
}
func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	result, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.Id = result.InsertedID.(primitive.ObjectID)
	if err := s.UpdateHotelRooms(ctx, room.HotelId, room.Id); err != nil {
		return nil, err
	}
	return room, nil
}
