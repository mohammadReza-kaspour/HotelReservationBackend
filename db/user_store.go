package db

import (
	"context"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

type UserStore interface {
	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
}

type MongoUserStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client:     client,
		collection: client.Database(dbName).Collection(userCollection),
	}
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.User
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	var users []*types.User
	cursor, err := s.collection.Find(ctx, bson.M{})
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
func(s *MongoUserStore) InsertUser(ctx context.Context ,user *types.User) (*types.User , error){
	result , err := s.collection.InsertOne(ctx , user)
	if err != nil {
		return nil , err
	}
	user.Id = result.InsertedID.(primitive.ObjectID)
	return user , nil
}