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

type UserHandler struct {
	storage db.Storage
}

func NewUserHandler(storage db.Storage) *UserHandler {
	return &UserHandler{
		storage: storage,
	}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("user").(*jwt.CustomClaims)
	if !ok {
		return localErrors.NewErrorUnAuthorized()
	}

	user, err := h.storage.UserStorage.GetUserById(c.Context(), userInfo.UserID)
	if err != nil {
		return err
	}

	if user == nil {
		return localErrors.NewErrorNotFound("user not found")
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", user)
	return c.Status(response.Code).JSON(response)
}
func (h *UserHandler) EditUser(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("user").(*jwt.CustomClaims)
	if !ok {
		return localErrors.NewErrorUnAuthorized()
	}

	var params *types.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
		return err
	}
	if err := h.storage.UserStorage.UpdateUser(c.Context(), userInfo.UserID, params); err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", &map[string]any{})
	return c.Status(response.Code).JSON(response)
}
