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

type UserStorage interface {
	shared.Dropper

	GetUserById(ctx context.Context, id string) (*types.User, error)
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	GetUsers(ctx context.Context, size *int64, page *int64) ([]*types.User, error)
	InsertUser(ctx context.Context, user *types.User) (*types.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, params *types.UpdateUserParams) error
}

type MongoUserStorage struct {
	collection *mongo.Collection
}

func NewMongoUserStorage(db *mongo.Database) *MongoUserStorage {
	return &MongoUserStorage{
		collection: db.Collection(config.DbUserCollectionName),
	}
}
func (s *MongoUserStorage) Drop(ctx context.Context) error {
	return s.collection.Drop(ctx)
}
func (s *MongoUserStorage) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user *types.User
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, localErrors.NewErrorNotFound("document not found")
		}
		return nil, err
	}

	return user, nil
}
func (s *MongoUserStorage) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
func (s *MongoUserStorage) GetUsers(ctx context.Context, size *int64, page *int64) ([]*types.User, error) {
	var users []*types.User
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
		var result *types.User
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		users = append(users, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
func (s *MongoUserStorage) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	result, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return user, nil
}
func (s *MongoUserStorage) DeleteUser(ctx context.Context, id string) error {
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
func (s *MongoUserStorage) UpdateUser(ctx context.Context, id string, params *types.UpdateUserParams) error {
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
