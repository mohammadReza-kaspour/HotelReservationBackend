package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomType int

const (
	_ RoomType = iota
	SingleRoomType
	DoubleRoomType
	SeaSideRoomType
	DeluxeRoomType
)

type Room struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType           `bson:"type" json:"type"`
	Price     float64            `bson:"price" json:"price"`
	HotelId   primitive.ObjectID `bson:"hotelId" json:"hotelId"`
}

func NewRoom(params *CreateRoomParams) *Room {
	return &Room{
		Type:      params.Type,
		Price:     params.Price,
		HotelId:   params.HotelId,
	}
}

type CreateRoomParams struct {
	Type      RoomType           `bson:"type" json:"type"`
	Price     float64            `bson:"price" json:"price"`
	HotelId   primitive.ObjectID `bson:"hotelId" json:"hotelId"`
}
