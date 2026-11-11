package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/jwt"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
)

type UserHandler struct {
	userStorage db.UserStorage
}

func NewUserHandler(userStorage db.UserStorage) *UserHandler {
	return &UserHandler{
		userStorage: userStorage,
	}
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStorage.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}
func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("user").(*jwt.CustomClaims)
	if !ok {
		return errors.New("userInfo not found")
	}

	user, err := h.userStorage.GetUserById(c.Context(), userInfo.UserID)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
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

	insertedUser, err := h.userStorage.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(insertedUser)
}
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	if err := h.userStorage.DeleteUser(c.Context(), userId); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": userId})
}
func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		userId = c.Params("id")
		params *types.UpdateUserParams
	)

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	//input validation
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	if err := h.userStorage.UpdateUser(c.Context(), userId, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": userId})
}
