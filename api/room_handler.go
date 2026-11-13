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

type RoomHandler struct {
	storage db.Storage
}

func NewRoomHandler(storage db.Storage) *RoomHandler {
	return &RoomHandler{
		storage: storage,
	}
}

func (h *RoomHandler) AddRoom(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("user").(*jwt.CustomClaims)
	if !ok {
		return localErrors.NewErrorUnAuthorized()
	}

	var params *types.CreateRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
		return err
	}

	//prove that user is the owner of hotel
	_, err := h.storage.HotelStorage.GetHotelOfOwnerById(c.Context(), userInfo.UserID, params.HotelId)
	if err != nil {
		return err
	}

	room, err := types.NewRoom(params)
	if err != nil {
		return err
	}
	fetchedRoom, err := h.storage.RoomStorage.InsertRoom(c.Context(), room)
	if err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", fetchedRoom)
	return c.Status(response.Code).JSON(response)
}
func (h *RoomHandler) GetRoom(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("user").(*jwt.CustomClaims)
	if !ok {
		return localErrors.NewErrorUnAuthorized()
	}

	var getRoomParams *utils.GetRoomParams
	if err := c.BodyParser(&getRoomParams); err != nil {
		return err
	}

	if err := getRoomParams.Validate(); err != nil {
		return err
	}

	//prove that user is the owner of hotel
	_, err := h.storage.HotelStorage.GetHotelOfOwnerById(c.Context(), userInfo.UserID, getRoomParams.HotelId)
	if err != nil {
		return err
	}

	fetchedRoom, err := h.storage.RoomStorage.GetRoomById(c.Context(), getRoomParams.RoomId)
	if err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", fetchedRoom)
	return c.Status(response.Code).JSON(response)
}
func (h *RoomHandler) GetRooms(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("user").(*jwt.CustomClaims)
	if !ok {
		return localErrors.NewErrorUnAuthorized()
	}

	var getRoomParams *utils.GetRoomPagination
	if err := c.BodyParser(&getRoomParams); err != nil {
		return err
	}
	if err := getRoomParams.Validate(); err != nil {
		return err
	}

	//prove that user is the owner of hotel
	_, err := h.storage.HotelStorage.GetHotelOfOwnerById(c.Context(), userInfo.UserID, getRoomParams.HotelId)
	if err != nil {
		return err
	}

	fetchedRooms, err := h.storage.RoomStorage.GetRoomsOfHotel(c.Context(), getRoomParams.HotelId, getRoomParams.Size, getRoomParams.Page)
	if err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", &fetchedRooms)
	return c.Status(response.Code).JSON(response)
}
func (h *RoomHandler) GetRoomPublic(c *fiber.Ctx) error {
	var getRoomParams *utils.GetId
	if err := c.BodyParser(&getRoomParams); err != nil {
		return err
	}

	if err := getRoomParams.Validate(); err != nil {
		return err
	}

	fetchedRoom, err := h.storage.RoomStorage.GetRoomById(c.Context(), getRoomParams.Id)
	if err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", fetchedRoom)
	return c.Status(response.Code).JSON(response)
}
func (h *RoomHandler) GetRoomsPublic(c *fiber.Ctx) error {
	var getRoomParams *utils.GetPagination
	if err := c.BodyParser(&getRoomParams); err != nil {
		return err
	}
	if err := getRoomParams.Validate(); err != nil {
		return err
	}

	fetchedRooms, err := h.storage.RoomStorage.GetRooms(c.Context(), getRoomParams.Size, getRoomParams.Page)
	if err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", &fetchedRooms)
	return c.Status(response.Code).JSON(response)
}
func (h *RoomHandler) DeleteRoom(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("user").(*jwt.CustomClaims)
	if !ok {
		return localErrors.NewErrorUnAuthorized()
	}

	var getRoomParams *utils.GetRoomParams
	if err := c.BodyParser(&getRoomParams); err != nil {
		return err
	}

	if err := getRoomParams.Validate(); err != nil {
		return err
	}

	//prove that user is the owner of hotel
	_, err := h.storage.HotelStorage.GetHotelOfOwnerById(c.Context(), userInfo.UserID, getRoomParams.HotelId)
	if err != nil {
		return err
	}

	if err := h.storage.RoomStorage.DeleteRoom(c.Context(), getRoomParams.RoomId); err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", &map[string]any{})
	return c.Status(response.Code).JSON(response)
}
func (h *RoomHandler) UpdateRoom(c *fiber.Ctx) error {
	userInfo, ok := c.Locals("user").(*jwt.CustomClaims)
	if !ok {
		return localErrors.NewErrorUnAuthorized()
	}

	var updateRoomParams *types.UpdateRoomParams
	if err := c.BodyParser(&updateRoomParams); err != nil {
		return err
	}
	if err := updateRoomParams.Validate(); err != nil {
		return err
	}

	//prove that user is the owner of hotel
	_, err := h.storage.HotelStorage.GetHotelOfOwnerById(c.Context(), userInfo.UserID, updateRoomParams.HotelId)
	if err != nil {
		return err
	}

	if err := h.storage.RoomStorage.UpdateRoom(c.Context(), updateRoomParams.Id, updateRoomParams); err != nil {
		return err
	}

	response := localResponse.NewResponse(http.StatusOK, "OK", &map[string]any{})
	return c.Status(response.Code).JSON(response)
}
