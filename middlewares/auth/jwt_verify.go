package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/jwt"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
)

type AuthHeaders struct {
	Authorization string
}

func VerifyJwtTokenMiddleware(c *fiber.Ctx) error {
	var headers AuthHeaders
	if err := c.ReqHeaderParser(&headers); err != nil {
		return err
	}

	jwtClaim, err := jwt.VerifyJwtToken(headers.Authorization)
	if err != nil {
		return localErrors.NewErrorUnAuthorized(err.Error())
	}

	c.Locals("user", jwtClaim)
	return c.Next()
}
