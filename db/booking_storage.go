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

type BookingStorage interface {
	shared.Dropper

	GetBookingById(ctx context.Context, id string) (*types.Booking, error)
	GetBookings(ctx context.Context, size *int64, page *int64) ([]*types.Booking, error)
	InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error)
	DeleteBooking(ctx context.Context, id string) error
	UpdateBooking(ctx context.Context, id string, params *types.UpdateBookingParams) error
}

type MongoBookingStorage struct {
	collection *mongo.Collection
}

func NewMongoBookingStorage(db *mongo.Database) *MongoBookingStorage {
	return &MongoBookingStorage{
		collection: db.Collection(config.DbBookingCollectionName),
	}
}
func (s *MongoBookingStorage) Drop(ctx context.Context) error {
	return s.collection.Drop(ctx)
}
func (s *MongoBookingStorage) GetBookingById(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var booking *types.Booking
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, localErrors.NewErrorNotFound("document not found")
		}
		return nil, err
	}

	return booking, nil
}
func (s *MongoBookingStorage) GetBookings(ctx context.Context, size *int64, page *int64) ([]*types.Booking, error) {
	var bookings []*types.Booking
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
		var result *types.Booking
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		bookings = append(bookings, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}
func (s *MongoBookingStorage) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	result, err := s.collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return booking, nil
}
func (s *MongoBookingStorage) DeleteBooking(ctx context.Context, id string) error {
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
func (s *MongoBookingStorage) UpdateBooking(ctx context.Context, id string, params *types.UpdateBookingParams) error {
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
