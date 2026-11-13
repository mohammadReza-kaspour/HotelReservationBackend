package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
)

type AuthHeaders struct {
	Authorization string
}

func (params *AuthHeaders) HashToken() string {
	sum := sha256.Sum256([]byte(params.Authorization))
	return hex.EncodeToString(sum[:])
}

type RefreshTokenParams struct {
	RefreshToken string `json:"refreshToken"`
}

func (params *RefreshTokenParams) Validate() error {
	if len(params.RefreshToken) == 0 {
		return localErrors.NewErrorBadRequest("refreshToken is required")
	}

	return nil
}
func (params *RefreshTokenParams) HashToken() string {
	sum := sha256.Sum256([]byte(params.RefreshToken))
	return hex.EncodeToString(sum[:])
}

type GetPagination struct {
	Size *int64 `json:"size"`
	Page *int64 `json:"page"`
}

func (params *GetPagination) Validate() error {
	if params.Size != nil && *params.Size < 0 {
		return localErrors.NewErrorBadRequest("size must be greater that 0")
	}
	if params.Page != nil && *params.Page < 1 {
		return localErrors.NewErrorBadRequest("page must be greater that 1")
	}
	return nil
}

type GetId struct {
	Id string `json:"id"`
}

func (params *GetId) Validate() error {
	if len(strings.TrimSpace(params.Id)) == 0 {
		return localErrors.NewErrorBadRequest("Id is required")
	}
	return nil
}

type GetRoomParams struct {
	HotelId string `json:"hotelId"`
	RoomId  string `json:"roomId"`
}

func (params *GetRoomParams) Validate() error {
	if len(strings.TrimSpace(params.HotelId)) == 0 {
		return localErrors.NewErrorBadRequest("hotelId is required")
	}
	if len(strings.TrimSpace(params.RoomId)) == 0 {
		return localErrors.NewErrorBadRequest("roomId is required")
	}
	return nil
}

type GetRoomPagination struct {
	HotelId string `json:"hotelId"`
	Size    *int64 `json:"size"`
	Page    *int64 `json:"page"`
}

func (params *GetRoomPagination) Validate() error {
	if len(strings.TrimSpace(params.HotelId)) == 0 {
		return localErrors.NewErrorBadRequest("hotelId is required")
	}
	if params.Size != nil && *params.Size < 0 {
		return localErrors.NewErrorBadRequest("size must be greater that 0")
	}
	if params.Page != nil && *params.Page < 1 {
		return localErrors.NewErrorBadRequest("page must be greater that 1")
	}
	return nil
}
