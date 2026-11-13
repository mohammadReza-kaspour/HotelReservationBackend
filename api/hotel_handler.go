package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/jwt"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localResponse"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types/utils"
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
func (h *HotelHandler) GetHotel(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("user").(*jwt.CustomClaims)
	if !ok {
		return localErrors.NewErrorUnAuthorized()
	}

	var getHotelParams *utils.GetId
	if err := c.BodyParser(&getHotelParams); err != nil {
		return err
	}

	if err := getHotelParams.Validate(); err != nil {
		return err
	}

	fetchedHotel, err := h.storage.HotelStorage.GetHotelOfOwnerById(c.Context(), userInfo.UserID, getHotelParams.Id)
	if err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", fetchedHotel)
	return c.Status(response.Code).JSON(response)
}
func (h *HotelHandler) GetHotels(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("user").(*jwt.CustomClaims)
	if !ok {
		return localErrors.NewErrorUnAuthorized()
	}

	var getHotelParams *utils.GetPagination
	if err := c.BodyParser(&getHotelParams); err != nil {
		return err
	}
	if err := getHotelParams.Validate(); err != nil {
		return err
	}

	fetchedHotels, err := h.storage.HotelStorage.GetHotelsOfOwner(c.Context(), userInfo.UserID, getHotelParams.Size, getHotelParams.Page)
	if err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", &fetchedHotels)
	return c.Status(response.Code).JSON(response)
}
func (h *HotelHandler) GetHotelPublic(c *fiber.Ctx) error {
	var getHotelParams *utils.GetId
	if err := c.BodyParser(&getHotelParams); err != nil {
		return err
	}

	if err := getHotelParams.Validate(); err != nil {
		return err
	}

	fetchedHotel, err := h.storage.HotelStorage.GetHotelById(c.Context(), getHotelParams.Id)
	if err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", fetchedHotel)
	return c.Status(response.Code).JSON(response)
}
func (h *HotelHandler) GetHotelsPublic(c *fiber.Ctx) error {
	var getHotelParams *utils.GetPagination
	if err := c.BodyParser(&getHotelParams); err != nil {
		return err
	}
	if err := getHotelParams.Validate(); err != nil {
		return err
	}

	fetchedHotels, err := h.storage.HotelStorage.GetHotels(c.Context(), getHotelParams.Size, getHotelParams.Page)
	if err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", &fetchedHotels)
	return c.Status(response.Code).JSON(response)
}
func (h *HotelHandler) DeleteHotel(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("user").(*jwt.CustomClaims)
	if !ok {
		return localErrors.NewErrorUnAuthorized()
	}

	var getHotelParams *utils.GetId
	if err := c.BodyParser(&getHotelParams); err != nil {
		return err
	}

	if err := getHotelParams.Validate(); err != nil {
		return err
	}

	if err := h.storage.HotelStorage.DeleteHotelOfOwner(c.Context(), userInfo.UserID, getHotelParams.Id); err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", &map[string]any{})
	return c.Status(response.Code).JSON(response)
}
func (h *HotelHandler) UpdateHotel(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("user").(*jwt.CustomClaims)
	if !ok {
		return localErrors.NewErrorUnAuthorized()
	}

	var updateHotelParams *types.UpdateHotelParams
	if err := c.BodyParser(&updateHotelParams); err != nil {
		return err
	}
	if err := updateHotelParams.Validate(); err != nil {
		return err
	}

	if err := h.storage.HotelStorage.UpdateHotelOfOwner(c.Context(), userInfo.UserID, updateHotelParams.Id, updateHotelParams); err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", &map[string]any{})
	return c.Status(response.Code).JSON(response)
}
