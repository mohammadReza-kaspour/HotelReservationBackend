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

type RoomStorage interface {
	shared.Dropper

	GetRoomById(ctx context.Context, id string) (*types.Room, error)
	GetRooms(ctx context.Context, size *int64, page *int64) ([]*types.Room, error)
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
	DeleteRoom(ctx context.Context, id string) error
	UpdateRoom(ctx context.Context, id string, params *types.UpdateRoomParams) error
}

type MongoRoomStorage struct {
	collection *mongo.Collection
}

func NewMongoRoomStorage(db *mongo.Database) *MongoRoomStorage {
	return &MongoRoomStorage{
		collection: db.Collection(config.DbRoomCollectionName),
	}
}
func (s *MongoRoomStorage) Drop(ctx context.Context) error {
	return s.collection.Drop(ctx)
}
func (s *MongoRoomStorage) GetRoomById(ctx context.Context, id string) (*types.Room, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var room *types.Room
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&room); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, localErrors.NewErrorNotFound("document not found")
		}
		return nil, err
	}

	return room, nil
}
func (s *MongoRoomStorage) GetRooms(ctx context.Context, size *int64, page *int64) ([]*types.Room, error) {
	var rooms []*types.Room
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
		var result *types.Room
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		rooms = append(rooms, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}
func (s *MongoRoomStorage) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	result, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return room, nil
}
func (s *MongoRoomStorage) DeleteRoom(ctx context.Context, id string) error {
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
func (s *MongoRoomStorage) UpdateRoom(ctx context.Context, id string, params *types.UpdateRoomParams) error {
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
