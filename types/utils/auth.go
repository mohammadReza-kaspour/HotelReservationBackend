package utils

import (
	"crypto/sha256"
	"encoding/hex"

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
