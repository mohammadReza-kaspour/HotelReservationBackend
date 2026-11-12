package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/jwt"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localResponse"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
)

type HotelHandler struct {
	storage db.Storage
}

func NewHotelHandler(storage db.Storage) *HotelHandler {
	return &HotelHandler{
		storage: storage,
	}
}

func (h *HotelHandler) AddHotel(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("user").(*jwt.CustomClaims)
	if !ok {
		return localErrors.NewErrorUnAuthorized()
	}

	var params *types.CreateHotelParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	params.Owner = userInfo.UserID
	if err := params.Validate(); err != nil {
		return err
	}

	hotel, err := types.NewHotel(params)
	if err != nil {
		return err
	}
	fetchedHotel, err := h.storage.HotelStorage.InsertHotel(c.Context(), hotel)
	if err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", fetchedHotel)
	return c.Status(response.Code).JSON(response)
}
