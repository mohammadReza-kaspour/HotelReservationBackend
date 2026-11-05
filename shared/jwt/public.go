package jwt

import (
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/config"
)

func GenerateAccessToken(userID string, email string) (string, error) {
	return generateToken(userID, email, config.AccessTokenExpireDelay)
}

func GenerateRefreshToken(userID string, email string) (string, error) {
	return generateToken(userID, email, config.RefreshTokenExpireDelay)
}
func VerifyJwtToken(tokenString string) (*CustomClaims, error) {
	return validateToken(tokenString)
}
