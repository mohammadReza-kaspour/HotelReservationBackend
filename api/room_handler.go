package api

import "github.com/mohammadReza-kaspour/HotelReservationBackend/db"

type RoomHandler struct {
	roomStorage  db.RoomStorage
	hotelStorage db.HotelStorage
}

func NewRoomHandler(roomStorage db.RoomStorage, hotelStorage db.HotelStorage) *RoomHandler {
	return &RoomHandler{
		roomStorage:  roomStorage,
		hotelStorage: hotelStorage,
	}
}
