package auth

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/jwt"
)

type AuthHeaders struct {
	Authorization string
}

func VerifyJwtTokenMiddleware(c *fiber.Ctx) error {
	var headers AuthHeaders
	if err := c.ReqHeaderParser(&headers);err!=nil{
		return err
	}

	jwtClaim , err := jwt.VerifyJwtToken(headers.Authorization)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{"error": err.Error()})
	}
	
	c.Locals("user" , jwtClaim)
	return c.Next()
}