package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/config"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/jwt"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types/utils"
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
	// when user do signup always sessions of it is empty
	// so we dont have to do circular freshing sessions in db

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

	newSessionParams := &types.CreateSessionParams{
		UserId: insertedUser.Id,
		Token:  refreshToken,
	}
	h.storage.SessionStorage.InsertSession(c.Context(), types.NewSession(newSessionParams))

	return c.JSON(map[string]string{
		"AccessToken":  accessToken,
		"RefreshToken": refreshToken,
	})
}
func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	// when user do signin its possible to create sessions as much as possible
	// so we have to do circular freshing sessions in db to prevent infinit login possibility

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

	sessions, err := h.storage.SessionStorage.GetSessionsOfUser(c.Context(), fetchedUser.Id)
	if err != nil {
		return err
	}

	if len(sessions) >= config.MaxAvailableSessions {
		if err := h.storage.SessionStorage.DeleteSession(c.Context(), sessions[0].Id); err != nil {
			return err
		}
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

	newSessionParams := &types.CreateSessionParams{
		UserId: fetchedUser.Id,
		Token:  refreshToken,
	}
	h.storage.SessionStorage.InsertSession(c.Context(), types.NewSession(newSessionParams))

	return c.JSON(map[string]string{
		"AccessToken":  accessToken,
		"RefreshToken": refreshToken,
	})
}
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	//session of user must update and new token hash replace with older one

	var params *utils.RefreshTokenParams
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

	currentSession, err := h.storage.SessionStorage.GetSessionsOfUserByTokenHash(c.Context(), jwtClaims.UserID, params.HashToken())
	if err != nil {
		return err
	}
	if currentSession == nil {
		return localErrors.NewErrorUnAuthorized()
	}

	accessToken, err := jwt.GenerateAccessToken(jwtClaims.UserID, jwtClaims.Email)
	if err != nil {
		return err
	}

	refreshToken, err := jwt.GenerateRefreshToken(jwtClaims.UserID, jwtClaims.Email)
	if err != nil {
		return err
	}

	updateSessionParams := &types.UpdateSessionParams{
		Token: refreshToken,
	}
	if err := h.storage.SessionStorage.UpdateSession(c.Context(), currentSession.Id, updateSessionParams); err != nil {
		return err
	}

	return c.JSON(map[string]string{
		"AccessToken":  accessToken,
		"RefreshToken": refreshToken,
	})
}
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	//session of user must update and new token hash replace with older one

	var params *utils.RefreshTokenParams
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

	session, err := h.storage.SessionStorage.GetSessionsOfUserByTokenHash(c.Context(), jwtClaims.UserID, params.HashToken())
	if err != nil {
		return localErrors.NewErrorUnAuthorized()
	}
	if session == nil {
		return localErrors.NewErrorUnAuthorized()
	}

	if err := h.storage.SessionStorage.DeleteSession(c.Context(), session.Id); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(session)
}
