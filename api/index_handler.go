package api

import "github.com/mohammadReza-kaspour/HotelReservationBackend/db"

type Handler struct {
	AuthHandler *AuthHandler
}

func NewHandler(storage db.Storage) *Handler{
	return &Handler{
		AuthHandler: NewAuthHandler(storage),
	}
}