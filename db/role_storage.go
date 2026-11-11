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

type RoleStorage interface {
	shared.Dropper

	GetRoleById(ctx context.Context, id string) (*types.Role, error)
	GetRoles(ctx context.Context, size *int64, page *int64) ([]*types.Role, error)
	InsertRole(ctx context.Context, role *types.Role) (*types.Role, error)
	DeleteRole(ctx context.Context, id string) error
	UpdateRole(ctx context.Context, id string, params *types.UpdateRoleParams) error
}

type MongoRoleStorage struct {
	collection *mongo.Collection
}

func NewMongoRoleStorage(db *mongo.Database) *MongoRoleStorage {
	return &MongoRoleStorage{
		collection: db.Collection(config.DbRoleCollectionName),
	}
}
func (s *MongoRoleStorage) Drop(ctx context.Context) error {
	return s.collection.Drop(ctx)
}
func (s *MongoRoleStorage) GetRoleById(ctx context.Context, id string) (*types.Role, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var role *types.Role
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&role); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, localErrors.NewErrorNotFound("document not found")
		}
		return nil, err
	}

	return role, nil
}
func (s *MongoRoleStorage) GetRoles(ctx context.Context, size *int64, page *int64) ([]*types.Role, error) {
	var roles []*types.Role
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
		var result *types.Role
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		roles = append(roles, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}
func (s *MongoRoleStorage) InsertRole(ctx context.Context, role *types.Role) (*types.Role, error) {
	result, err := s.collection.InsertOne(ctx, role)
	if err != nil {
		return nil, err
	}
	role.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return role, nil
}
func (s *MongoRoleStorage) DeleteRole(ctx context.Context, id string) error {
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
func (s *MongoRoleStorage) UpdateRole(ctx context.Context, id string, params *types.UpdateRoleParams) error {
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
