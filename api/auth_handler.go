package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/jwt"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
)

type AuthHandler struct {
	storage db.Storage
}

func NewAuthHandler(storage db.Storage) *AuthHandler {
	return &AuthHandler{
		storage: storage,
	}
}

func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	var params *types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
		return err
	}

	user, err := types.NewUser(params)
	if err != nil {
		return err
	}

	fetchedUser, err := h.storage.UserStorage.GetUserByEmailOrNil(c.Context(), params.Email)
	if err != nil {
		return err
	}
	if fetchedUser != nil {
		return localErrors.NewErrorConflict("user already exists")
	}

	insertedUser, err := h.storage.UserStorage.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	accessToken, err := jwt.GenerateAccessToken(insertedUser.Id, insertedUser.Email)
	if err != nil {
		return err
	}

	refreshToken, err := jwt.GenerateRefreshToken(insertedUser.Id, insertedUser.Email)
	if err != nil {
		return err
	}

	return c.JSON(map[string]string{
		"AccessToken":  accessToken,
		"RefreshToken": refreshToken,
	})
}
func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	var params *types.LoginUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
		return err
	}

	fetchedUser, err := h.storage.UserStorage.GetUserByEmailOrNil(c.Context(), params.Email)
	if err != nil {
		return err
	}
	if fetchedUser == nil {
		return localErrors.NewErrorUnAuthorized("wrong user name or password")
	}

	if err := fetchedUser.VerifyPassword(params.Password); err != nil {
		return localErrors.NewErrorUnAuthorized("wrong user name or password")
	}

	accessToken, err := jwt.GenerateAccessToken(fetchedUser.Id, fetchedUser.Email)
	if err != nil {
		return err
	}

	refreshToken, err := jwt.GenerateRefreshToken(fetchedUser.Id, fetchedUser.Email)
	if err != nil {
		return err
	}

	return c.JSON(map[string]string{
		"AccessToken":  accessToken,
		"RefreshToken": refreshToken,
	})
}
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var params *types.RefreshTokenParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
		return err
	}

	jwtClaims, err := jwt.VerifyJwtToken(params.RefreshToken)
	if err != nil {
		return err
	}

	accessToken, err := jwt.GenerateAccessToken(jwtClaims.UserID, jwtClaims.Email)
	if err != nil {
		return err
	}

	refreshToken, err := jwt.GenerateRefreshToken(jwtClaims.UserID, jwtClaims.Email)
	if err != nil {
		return err
	}

	return c.JSON(map[string]string{
		"AccessToken":  accessToken,
		"RefreshToken": refreshToken,
	})
}
