package api

import "github.com/mohammadReza-kaspour/HotelReservationBackend/db"

type RoomHandler struct {
	storage db.Storage
}

func NewRoomHandler(storage db.Storage) *RoomHandler {
	return &RoomHandler{
		storage: storage,
	}
}
