package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/jwt"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
)

type AuthHandler struct {
	userStorage db.UserStorage
}

func NewAuthHandler(userStorage db.UserStorage) *AuthHandler {
	return &AuthHandler{
		userStorage: userStorage,
	}
}

func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	var params *types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	//input validation
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	fetchedUser, err := h.userStorage.GetUserByEmailOrNil(c.Context(), params.Email)
	if err != nil {
		return err
	}
	if fetchedUser != nil {
		return errors.New("user already exists")
	}

	insertedUser, err := h.userStorage.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	accessToken, err := jwt.GenerateAccessToken(insertedUser.Id.Hex(), insertedUser.Email)
	if err != nil {
		return err
	}

	refreshToken, err := jwt.GenerateRefreshToken(insertedUser.Id.Hex(), insertedUser.Email)
	if err != nil {
		return err
	}

	return c.JSON(map[string]string{
		"AccessToken":  accessToken,
		"RefreshToken": refreshToken,
	})
}

func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	return c.Status(http.StatusNotImplemented).SendString("Not Implemented Yet")
}
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	return c.Status(http.StatusNotImplemented).SendString("Not Implemented Yet")
}
