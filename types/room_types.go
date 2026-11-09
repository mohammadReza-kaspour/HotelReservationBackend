package types

import (
	"fmt"
	"regexp"
	"time"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types/shared"
)

type Room struct {
	Id          string           `bson:"_id,omitempty" json:"id,omitempty"`
	HotelId     string           `bson:"hotelId" json:"hotelId"`
	Title       string           `bson:"title" json:"title"`
	BedType     shared.BedType   `bson:"bedType" json:"bedType"`
	Capacity    int              `bson:"capacity" json:"capacity"`
	SizeSqm     float64          `bson:"sizeSqm" json:"sizeSqm"`
	View        shared.RoomView  `bson:"view" json:"view"`
	Amenities   []shared.Amenity `bson:"amenities,omitempty" json:"amenities,omitempty"`
	Price       float64          `bson:"price" json:"price"`
	Currency    string           `bson:"currency" json:"currency"`
	IsAvailable bool             `bson:"isAvailable" json:"isAvailable"`
	Version     int              `bson:"version" json:"version"` //Optimistic lock guard for concurrent booking
	CreatedAt   time.Time        `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time        `bson:"updatedAt" json:"updatedAt"`
}

func NewRoom(params *CreateRoomParams) *Room {
	return &Room{
		HotelId:     params.HotelId,
		Title:       params.Title,
		BedType:     params.BedType,
		Capacity:    params.Capacity,
		SizeSqm:     params.SizeSqm,
		View:        params.View,
		Amenities:   params.Amenities,
		Price:       params.Price,
		Currency:    "IR",
		IsAvailable: true,
		Version:     0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

type CreateRoomParams struct {
	HotelId   string           `json:"hotelId"`
	Title     string           `json:"title"`
	BedType   shared.BedType   `json:"bedType"`
	Capacity  int              `json:"capacity"`
	SizeSqm   float64          `json:"sizeSqm"`
	View      shared.RoomView  `json:"view"`
	Amenities []shared.Amenity `json:"amenities,omitempty"`
	Price     float64          `json:"price"`
}

func (params *CreateRoomParams) Validate() error {
	validStringRegexp := regexp.MustCompile(`^[a-zA-Z0-9_,\-]+$`)

	if len(params.HotelId) == 0 {
		return localErrors.NewErrorBadRequest("room hotelId is required")
	}
	if !validStringRegexp.MatchString(params.Title) {
		return localErrors.NewErrorBadRequest("room title is invalid")
	}
	if !params.BedType.Validate() {
		return localErrors.NewErrorBadRequest("room bedType is invalid")
	}
	if !params.View.Validate() {
		return localErrors.NewErrorBadRequest("room view is invalid")
	}

	if params.Capacity < 0 {
		return localErrors.NewErrorBadRequest("room capacity can't be negative")
	}
	if params.SizeSqm < 0 {
		return localErrors.NewErrorBadRequest("room sizeSqm can't be negative")
	}
	if params.Price < 0 {
		return localErrors.NewErrorBadRequest("room price can't be negative")
	}

	for _, amenity := range params.Amenities {
		if !amenity.Validate() {
			return localErrors.NewErrorBadRequest(fmt.Sprintf("amenity %s is invalid", amenity))
		}
	}

	return nil
}

type UpdateRoomParams struct {
	HotelId     string           `bson:"hotelId" json:"hotelId"`
	Title       string           `bson:"title" json:"title"`
	BedType     shared.BedType   `bson:"bedType" json:"bedType"`
	Capacity    int              `bson:"capacity" json:"capacity"`
	SizeSqm     float64          `bson:"sizeSqm" json:"sizeSqm"`
	View        shared.RoomView  `bson:"view" json:"view"`
	Amenities   []shared.Amenity `bson:"amenities,omitempty" json:"amenities,omitempty"`
	Price       float64          `bson:"price" json:"price"`
	IsAvailable bool             `bson:"isAvailable" json:"isAvailable"`
}

func (params *UpdateRoomParams) Validate() error {
	validStringRegexp := regexp.MustCompile(`^[a-zA-Z0-9_,\-]+$`)

	if len(params.HotelId) == 0 {
		return localErrors.NewErrorBadRequest("room hotelId is required")
	}
	if params.Title != "" && !validStringRegexp.MatchString(params.Title) {
		return localErrors.NewErrorBadRequest("room title is invalid")
	}
	if params.BedType != "" && !params.BedType.Validate() {
		return localErrors.NewErrorBadRequest("room bedType is invalid")
	}
	if params.View != "" && !params.View.Validate() {
		return localErrors.NewErrorBadRequest("room view is invalid")
	}

	if params.Capacity != 0 && params.Capacity < 0 {
		return localErrors.NewErrorBadRequest("room capacity can't be negative")
	}
	if params.SizeSqm != 0 && params.SizeSqm < 0 {
		return localErrors.NewErrorBadRequest("room sizeSqm can't be negative")
	}
	if params.Price != 0 && params.Price < 0 {
		return localErrors.NewErrorBadRequest("room price can't be negative")
	}

	if params.Amenities != nil {
		for _, amenity := range params.Amenities {
			if !amenity.Validate() {
				return localErrors.NewErrorBadRequest(fmt.Sprintf("amenity %s is invalid", amenity))
			}
		}
	}

	return nil
}
