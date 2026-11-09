package types

import (
	"fmt"
	"regexp"
	"time"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types/shared"
)

type Hotel struct {
	Id        string           `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string           `bson:"name" json:"name"`
	City      string           `bson:"city" json:"city"`
	Country   string           `bson:"country" json:"country"`
	Address   string           `bson:"address" json:"address"`
	Rate      int              `bson:"rate" json:"rate"`
	Amenities []shared.Amenity `bson:"amenities,omitempty" json:"amenities,omitempty"`
	Rooms     []string         `bson:"rooms" json:"rooms"`
	CreatedAt time.Time        `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time        `bson:"updatedAt" json:"updatedAt"`
}

func NewHotel(params *CreateHotelParams) *Hotel {

	initAmenities := []shared.Amenity{}
	if params.Amenities != nil {
		initAmenities = params.Amenities
	}

	return &Hotel{
		Name:      params.Name,
		City:      params.City,
		Country:   params.Country,
		Address:   params.Address,
		Rate:      params.Rate,
		Amenities: initAmenities,
		Rooms:     []string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type CreateHotelParams struct {
	// creation doesnt need bson cause it will pass into type constructor
	Name      string           `json:"name"`
	City      string           `json:"city"`
	Country   string           `json:"country"`
	Address   string           `json:"address"`
	Rate      int              `json:"rate"`
	Amenities []shared.Amenity `json:"amenities,omitempty"`
}

func (params *CreateHotelParams) Validate() error {
	validShortStringRegexp := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	validLongStringRegexp := regexp.MustCompile(`^[a-zA-Z0-9 _,\-]+$`)

	if !validShortStringRegexp.MatchString(params.Name) {
		return localErrors.NewErrorBadRequest("hotel name is invalid")
	}
	if !validShortStringRegexp.MatchString(params.City) {
		return localErrors.NewErrorBadRequest("hotel city is invalid")
	}
	if !validShortStringRegexp.MatchString(params.Country) {
		return localErrors.NewErrorBadRequest("hotel country is invalid")
	}
	if !validLongStringRegexp.MatchString(params.Address) {
		return localErrors.NewErrorBadRequest("hotel address is invalid")
	}

	if params.Rate < 1 || params.Rate > 5 {
		return localErrors.NewErrorBadRequest("rate must be between 1 and 5")
	}

	for _, amenity := range params.Amenities {
		if !amenity.Validate() {
			return localErrors.NewErrorBadRequest(fmt.Sprintf("amenity %s is invalid", amenity))
		}
	}

	return nil
}

type UpdateHotelParams struct {
	// update need bson cause it will pass directly into mongodb driver
	Id        string           `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string           `bson:"name,omitempty" json:"name,omitempty"`
	City      string           `bson:"city,omitempty" json:"city,omitempty"`
	Country   string           `bson:"country,omitempty" json:"country,omitempty"`
	Address   string           `bson:"address,omitempty" json:"address,omitempty"`
	Rate      int              `bson:"rate,omitempty" json:"rate,omitempty"`
	Amenities []shared.Amenity `bson:"amenities,omitempty" json:"amenities,omitempty"`
}

func (params *UpdateHotelParams) Validate() error {
	validShortStringRegexp := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	validLongStringRegexp := regexp.MustCompile(`^[a-zA-Z0-9 _,\-]+$`)

	if len(params.Id) == 0 {
		return localErrors.NewErrorBadRequest("hotel id is required")
	}
	if params.Name != "" && !validShortStringRegexp.MatchString(params.Name) {
		return localErrors.NewErrorBadRequest("hotel name is invalid")
	}
	if params.City != "" && !validShortStringRegexp.MatchString(params.City) {
		return localErrors.NewErrorBadRequest("hotel city is invalid")
	}
	if params.Country != "" && !validShortStringRegexp.MatchString(params.Country) {
		return localErrors.NewErrorBadRequest("hotel country is invalid")
	}
	if params.Address != "" && !validLongStringRegexp.MatchString(params.Address) {
		return localErrors.NewErrorBadRequest("hotel address is invalid")
	}

	if params.Rate != 0 && params.Rate < 1 || params.Rate > 5 {
		return localErrors.NewErrorBadRequest("rate must be between 1 and 5")
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
