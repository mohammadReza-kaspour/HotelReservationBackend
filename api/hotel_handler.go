package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
)

type HotelHandler struct {
	hotelStorage db.HotelStorage
}

func NewHotelHandler(hotelStorage db.HotelStorage) *HotelHandler {
	return &HotelHandler{
		hotelStorage: hotelStorage,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.hotelStorage.GetHotels(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}
