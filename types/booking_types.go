package types

import (
	"time"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types/shared"
)

type Booking struct {
	Id            string               `bson:"_id,omitempty" json:"id,omitempty"`
	UserId        string               `bson:"userId" json:"userId"`
	RoomId        string               `bson:"roomId" json:"roomId"`
	CheckIn       time.Time            `bson:"checkIn" json:"checkIn"`
	CheckOut      time.Time            `bson:"checkOut" json:"checkOut"`
	Guests        int                  `bson:"guests" json:"guests"`
	TotalPrice    float64              `bson:"totalPrice" json:"totalPrice"`
	Currency      string               `bson:"currency" json:"currency"`
	Status        shared.BookingStatus `bson:"status" json:"status"`
	PaymentStatus shared.PaymentStatus `bson:"paymentStatus" json:"paymentStatus"`
	Version       int                  `bson:"version" json:"version"` //Optimistic lock guard for concurrent booking
	CreatedAt     time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time            `bson:"updatedAt" json:"updatedAt"`
}

func NewBooking(params *CreateBookingParams) *Booking {
	now := time.Now()
	return &Booking{
		UserId:        params.UserId,
		RoomId:        params.RoomId,
		CheckIn:       params.CheckIn,
		CheckOut:      params.CheckOut,
		Guests:        params.Guests,
		TotalPrice:    params.TotalPrice,
		Currency:      params.Currency,
		Status:        shared.BookingPending,
		PaymentStatus: shared.PaymentUnpaid,
		Version:       0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

type CreateBookingParams struct {
	UserId     string    `json:"userId"`
	RoomId     string    `json:"roomId"`
	CheckIn    time.Time `json:"checkIn"`
	CheckOut   time.Time `json:"checkOut"`

	Guests     int       `json:"guests"`     //must inherit from room detail
	TotalPrice float64   `json:"totalPrice"` //must inherit from room detail
	Currency   string    `json:"currency"`   //must inherit from room detail
}

func (params *CreateBookingParams) Validate() error {
	now := time.Now()

	if len(params.UserId) == 0 {
		return localErrors.NewErrorBadRequest("booking userId is required")
	}
	if len(params.RoomId) == 0 {
		return localErrors.NewErrorBadRequest("booking roomId is required")
	}

	if params.Guests <= 0 {
		return localErrors.NewErrorBadRequest("booking guests is invalid")
	}
	if params.TotalPrice <= 0 {
		return localErrors.NewErrorBadRequest("booking totalPrice is invalid")
	}
	if len(params.Currency) == 0 {
		return localErrors.NewErrorBadRequest("booking currency is invalid")
	}

	if params.CheckIn.Before(now) || params.CheckOut.Before(now) || params.CheckOut.Before(params.CheckIn) {
		return localErrors.NewErrorBadRequest("booking checkIn/checkOut is invalid")
	}

	return nil
}

type UpdateBookingParams struct {
	UserId   string     `bson:"userId" json:"userId"`
	RoomId   string     `bson:"roomId" json:"roomId"`
	CheckIn  *time.Time `bson:"checkIn" json:"checkIn"`
	CheckOut *time.Time `bson:"checkOut" json:"checkOut"`
}

func (params *UpdateBookingParams) Validate() error {
	now := time.Now()

	if len(params.UserId) == 0 {
		return localErrors.NewErrorBadRequest("booking userId is required")
	}
	if len(params.RoomId) == 0 {
		return localErrors.NewErrorBadRequest("booking roomId is required")
	}

	if params.CheckIn != nil && params.CheckIn.Before(now) {
		return localErrors.NewErrorBadRequest("booking checkIn/checkOut is invalid")
	}
	if params.CheckOut != nil && params.CheckOut.Before(now) {
		return localErrors.NewErrorBadRequest("booking checkIn/checkOut is invalid")
	}
	if params.CheckIn != nil && params.CheckOut != nil && params.CheckOut.Before(*params.CheckIn) {
		return localErrors.NewErrorBadRequest("booking checkIn/checkOut is invalid")
	}

	return nil
}

func (params *UpdateBookingParams) ToUpdateMap() map[string]any {
	updateMap := map[string]any{}

	updateMap["userId"] = params.UserId
	updateMap["roomId"] = params.RoomId

	derefNullable(updateMap , "checkIn" , params.CheckIn)
	derefNullable(updateMap , "checkOut" , params.CheckOut)

	if len(updateMap) == 0 {
		return nil
	}
	return updateMap
}