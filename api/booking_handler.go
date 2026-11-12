package api

import "github.com/mohammadReza-kaspour/HotelReservationBackend/db"

type BookingHandler struct {
	storage db.Storage
}

func NewBookingHandler(storage db.Storage) *BookingHandler {
	return &BookingHandler{
		storage: storage,
	}
}


