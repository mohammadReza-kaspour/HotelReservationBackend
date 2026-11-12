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

type SessionStorage interface {
	shared.Dropper
	shared.Statistician

	GetSessionById(ctx context.Context, id string) (*types.Session, error)
	GetSessionsOfUser(ctx context.Context, userId string) ([]*types.Session, error)
	GetSessionsOfUserByTokenHash(ctx context.Context, userId string, tokenHash string) (*types.Session, error)
	InsertSession(ctx context.Context, session *types.Session) (*types.Session, error)
	DeleteSession(ctx context.Context, id string) error
	UpdateSession(ctx context.Context, id string, params *types.UpdateSessionParams) error
}

type MongoSessionStorage struct {
	collection *mongo.Collection
}

func NewMongoSessionStorage(db *mongo.Database) *MongoSessionStorage {
	return &MongoSessionStorage{
		collection: db.Collection(config.DbSessionCollectionName),
	}
}
func (s *MongoSessionStorage) Drop(ctx context.Context) error {
	return s.collection.Drop(ctx)
}
func (s *MongoSessionStorage) Info(ctx context.Context) (*shared.StorageInfoSchema, error) {
	documentCount, err := s.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	info := shared.StorageInfoSchema{
		Counts: int(documentCount),
	}
	return &info, nil
}
func (s *MongoSessionStorage) GetSessionById(ctx context.Context, id string) (*types.Session, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var session *types.Session
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&session); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, localErrors.NewErrorNotFound("document not found")
		}
		return nil, err
	}

	return session, nil
}
func (s *MongoSessionStorage) GetSessionsOfUser(ctx context.Context, userId string) ([]*types.Session, error) {
	var sessions []*types.Session

	cursor, err := s.collection.Find(
		ctx,
		bson.M{"userId": userId},
		&options.FindOptions{
			Sort: bson.D{{Key: "createdAt", Value: 1}},
		},
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result *types.Session
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		sessions = append(sessions, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}
func (s *MongoSessionStorage) GetSessionsOfUserByTokenHash(ctx context.Context, userId string, tokenHash string) (*types.Session, error) {
	var session *types.Session
	if err := s.collection.FindOne(ctx, bson.M{"userId": userId, "tokenHash": tokenHash}).Decode(&session); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return session, nil
}
func (s *MongoSessionStorage) InsertSession(ctx context.Context, session *types.Session) (*types.Session, error) {
	result, err := s.collection.InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}
	session.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return session, nil
}
func (s *MongoSessionStorage) DeleteSession(ctx context.Context, id string) error {
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
func (s *MongoSessionStorage) UpdateSession(ctx context.Context, id string, params *types.UpdateSessionParams) error {
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
