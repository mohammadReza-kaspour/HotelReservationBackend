package api

import "github.com/mohammadReza-kaspour/HotelReservationBackend/db"

type RoomHandler struct {
	roomStore db.RoomStore
	hotelStore db.HotelStore
}

func NewRoomHandler(roomStore db.RoomStore , hotelStore db.HotelStore) *RoomHandler {
	return &RoomHandler{
		roomStore: roomStore,
		hotelStore: hotelStore,
	}
}

