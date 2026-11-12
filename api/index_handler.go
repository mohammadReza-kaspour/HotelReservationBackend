package api

import "github.com/mohammadReza-kaspour/HotelReservationBackend/db"

type Handler struct {
	AuthHandler    *AuthHandler
	UserHandler    *UserHandler
	HotelHandler   *HotelHandler
	RoomHandler    *RoomHandler
	BookingHandler *BookingHandler
}

func NewHandler(storage db.Storage) *Handler {
	return &Handler{
		AuthHandler:    NewAuthHandler(storage),
		UserHandler:    NewUserHandler(storage),
		HotelHandler:   NewHotelHandler(storage),
		RoomHandler:    NewRoomHandler(storage),
		BookingHandler: NewBookingHandler(storage),
	}
}
