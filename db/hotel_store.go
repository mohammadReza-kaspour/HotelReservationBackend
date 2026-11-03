package db

import (
	"context"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	Dropper

	InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	UpdateHotel(ctx context.Context, id string, params *types.UpdateHotelParams) error
	UpdateHotelRooms(ctx context.Context, oid primitive.ObjectID, roomOid primitive.ObjectID) error
	GetHotels(ctx context.Context) ([]*types.Hotel, error)
}

type MongoHotelStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbName string) *MongoHotelStore {
	return &MongoHotelStore{
		client:     client,
		collection: client.Database(dbName).Collection(hotelCollection),
	}
}

func (s *MongoHotelStore) Drop(ctx context.Context) error {
	return s.collection.Drop(ctx)
}
func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	result, err := s.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}

	hotel.Id = result.InsertedID.(primitive.ObjectID)
	return hotel, nil
}
func (s *MongoHotelStore) UpdateHotel(ctx context.Context, id string, params *types.UpdateHotelParams) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}
	update := bson.D{
		{Key: "$set", Value: params},
	}

	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
func (s *MongoHotelStore) UpdateHotelRooms(ctx context.Context, oid primitive.ObjectID, roomOid primitive.ObjectID) error {
	filter := bson.D{
		{Key: "_id", Value: oid},
	}
	update := bson.D{
		{Key: "$push", Value: bson.M{"rooms": roomOid}},
	}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
func (s *MongoHotelStore) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	var hotels []*types.Hotel
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result *types.Hotel
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
		hotels = append(hotels, result)
	}
	return hotels, nil
}
