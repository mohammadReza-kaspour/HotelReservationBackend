package main

import (
	"flag"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/api"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
)

const dbName = "hotelReservation"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		// Check if it's a Api error
		if apiError, ok := err.(*localErrors.ErrorBase); ok {
			return c.Status(apiError.Code).JSON(apiError)
		}

		// Check if it's a Fiber error
		if fiberError, ok := err.(*fiber.Error); ok {
			apiError := localErrors.NewErrorCustom(fiberError.Code, fiberError.Message, map[string]any{})
			return c.Status(apiError.Code).JSON(apiError)
		}

		apiError := localErrors.NewErrorCustom(http.StatusInternalServerError, err.Error(), map[string]any{})
		return c.Status(apiError.Code).JSON(apiError)
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "the listen addr of api server")
	flag.Parse()

	storage := db.NewStorage()
	handlers := api.NewHandler(*storage)

	app := fiber.New(config)
	app.Get("/developer", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("ok")
	})

	app.Post("/signup", handlers.AuthHandler.SignUp)
	app.Post("/signin", handlers.AuthHandler.SignIn)
	app.Post("/refresh-token", handlers.AuthHandler.RefreshToken)
	app.Post("/logout", handlers.AuthHandler.Logout)

	app.Listen(*listenAddr)
}
