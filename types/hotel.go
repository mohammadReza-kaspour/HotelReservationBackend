package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	Id       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rate     int                  `bson:"rate" json:"rate"`
}

func NewHotel(params *CreateHotelParams) *Hotel {
	return &Hotel{
		Name:     params.Name,
		Location: params.Location,
		Rate: params.Rate,
		Rooms:    []primitive.ObjectID{},
	}
}

type CreateHotelParams struct {
	Name     string `bson:"name" json:"name"`
	Location string `bson:"location" json:"location"`
	Rate     int `bson:"rate" json:"rate"`
}

type UpdateHotelParams struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Location string             `bson:"location" json:"location"`
}
