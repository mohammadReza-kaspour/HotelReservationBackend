package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	Id       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
}

func NewHotel(parms *CreateHotelParams) *Hotel {
	return &Hotel{
		Name:     parms.Name,
		Location: parms.Location,
		Rooms:    []primitive.ObjectID{},
	}
}

type CreateHotelParams struct {
	Name     string `bson:"name" json:"name"`
	Location string `bson:"location" json:"location"`
}

type UpdateHotelParams struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Location string             `bson:"location" json:"location"`
}
