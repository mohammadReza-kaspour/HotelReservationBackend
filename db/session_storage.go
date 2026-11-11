package db

import (
	"context"
	"errors"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/db/shared"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/config"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionStorage interface {
	shared.Dropper

	GetSessionById(ctx context.Context, id string) (*types.Session, error)
	GetSessions(ctx context.Context) ([]*types.Session, error)
	InsertSession(ctx context.Context, session *types.Session) (*types.Session, error)
	DeleteSession(ctx context.Context, id string) error
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
func (s *MongoSessionStorage) GetSessions(ctx context.Context) ([]*types.Session, error) {
	var sessions []*types.Session

	cursor, err := s.collection.Find(
		ctx,
		bson.M{},
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
