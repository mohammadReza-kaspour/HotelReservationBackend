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

func NewRoom(params *CreateRoomParams) (*Room, error) {
	now := time.Now()
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
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
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
	validStringRegexp := regexp.MustCompile(`^[a-zA-Z0-9_,\- ]+$`)

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
	Id          string           `bson:"_id" json:"id"`
	HotelId     string           `bson:"hotelId" json:"hotelId"`
	Title       *string          `bson:"title" json:"title"`
	BedType     *shared.BedType  `bson:"bedType" json:"bedType"`
	Capacity    *int             `bson:"capacity" json:"capacity"`
	SizeSqm     *float64         `bson:"sizeSqm" json:"sizeSqm"`
	View        *shared.RoomView `bson:"view" json:"view"`
	Amenities   []shared.Amenity `bson:"amenities,omitempty" json:"amenities,omitempty"`
	Price       *float64         `bson:"price" json:"price"`
	IsAvailable *bool            `bson:"isAvailable" json:"isAvailable"`
}

func (params *UpdateRoomParams) Validate() error {
	validStringRegexp := regexp.MustCompile(`^[a-zA-Z0-9_,\- ]+$`)

	if len(params.Id) == 0 {
		return localErrors.NewErrorBadRequest("room Id is required")
	}
	if len(params.HotelId) == 0 {
		return localErrors.NewErrorBadRequest("room hotelId is required")
	}
	if params.Title != nil && !validStringRegexp.MatchString(*params.Title) {
		return localErrors.NewErrorBadRequest("room title is invalid")
	}
	if params.BedType != nil && !params.BedType.Validate() {
		return localErrors.NewErrorBadRequest("room bedType is invalid")
	}
	if params.View != nil && !params.View.Validate() {
		return localErrors.NewErrorBadRequest("room view is invalid")
	}

	if params.Capacity != nil && *params.Capacity < 0 {
		return localErrors.NewErrorBadRequest("room capacity can't be negative")
	}
	if params.SizeSqm != nil && *params.SizeSqm < 0 {
		return localErrors.NewErrorBadRequest("room sizeSqm can't be negative")
	}
	if params.Price != nil && *params.Price < 0 {
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
func (params *UpdateRoomParams) ToUpdateMap() map[string]any {
	updateMap := map[string]any{}

	derefNullable(updateMap, "title", params.Title)
	derefNullable(updateMap, "bedType", params.BedType)
	derefNullable(updateMap, "capacity", params.Capacity)
	derefNullable(updateMap, "sizeSqm", params.SizeSqm)
	derefNullable(updateMap, "view", params.View)
	derefNullable(updateMap, "price", params.Price)
	derefNullable(updateMap, "isAvailable", params.IsAvailable)

	if params.Amenities != nil {
		updateMap["amenities"] = params.Amenities
	}

	if len(updateMap) == 0 {
		return nil
	}
	return updateMap
}
