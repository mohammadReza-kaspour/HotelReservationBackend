package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}
func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userStore.GetUserById(c.Context(), id)
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

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(insertedUser)
}
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userId); err != nil {
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
	
	if err := h.userStore.UpdateUser(c.Context() , userId , params); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": userId})
}
