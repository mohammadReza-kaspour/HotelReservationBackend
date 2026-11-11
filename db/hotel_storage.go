package db

import (
	"context"
	"errors"
	"time"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/db/shared"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/config"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelStorage interface {
	shared.Dropper

	GetHotelById(ctx context.Context, id string) (*types.Hotel, error)
	GetHotels(ctx context.Context, size *int64, page *int64) ([]*types.Hotel, error)
	InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	DeleteHotel(ctx context.Context, id string) error
	UpdateHotel(ctx context.Context, id string, params *types.UpdateHotelParams) error
}

type MongoHotelStorage struct {
	collection *mongo.Collection
}

func NewMongoHotelStorage(db *mongo.Database) *MongoHotelStorage {
	return &MongoHotelStorage{
		collection: db.Collection(config.DbHotelCollectionName),
	}
}
func (s *MongoHotelStorage) Drop(ctx context.Context) error {
	return s.collection.Drop(ctx)
}
func (s *MongoHotelStorage) GetHotelById(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var hotel *types.Hotel
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, localErrors.NewErrorNotFound("document not found")
		}
		return nil, err
	}

	return hotel, nil
}
func (s *MongoHotelStorage) GetHotels(ctx context.Context, size *int64, page *int64) ([]*types.Hotel, error) {
	var hotels []*types.Hotel
	skip := (*page - 1) * *size

	cursor, err := s.collection.Find(
		ctx,
		bson.M{},
		&options.FindOptions{
			Skip:  &skip,
			Limit: size,
		},
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result *types.Hotel
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		hotels = append(hotels, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return hotels, nil
}
func (s *MongoHotelStorage) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	result, err := s.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return hotel, nil
}
func (s *MongoHotelStorage) DeleteHotel(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}
func (s *MongoHotelStorage) UpdateHotel(ctx context.Context, id string, params *types.UpdateHotelParams) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateMap := params.ToUpdateMap()
	updateMap["updatedAt"] = time.Now()

	mongoFilter := bson.D{{Key: "_id", Value: oid}}
	mongoUpdate := bson.D{
		{Key: "$set", Value: updateMap},
	}

	_, err = s.collection.UpdateOne(ctx, mongoFilter, mongoUpdate)
	if err != nil {
		return err
	}

	return nil
}
