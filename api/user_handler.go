package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
)

type UserHandler struct {
	userStore db.UserStore
}
func NewUserHandler(userStore db.UserStore) *UserHandler{
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "mohammad",
		LastName:  "kaspour",
	}
	return c.JSON([]types.User{u})
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()

	user , err := h.userStore.GetUserById(ctx , id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
