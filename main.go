package main

import (
	"flag"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/api"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/middlewares"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
)

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

	apiV1 := app.Group("/api/v1")

	apiV1Users := apiV1.Group("/user", middlewares.VerifyJwtTokenMiddleware)
	apiV1Users.Get("/", handlers.UserHandler.GetUser)
	apiV1Users.Put("/", handlers.UserHandler.EditUser)

	apiV1Hotel := apiV1.Group("/hotel", middlewares.VerifyJwtTokenMiddleware)
	apiV1Hotel.Post("/add", handlers.HotelHandler.AddHotel)
	apiV1Hotel.Post("/list", handlers.HotelHandler.GetHotels)
	apiV1Hotel.Post("/get", handlers.HotelHandler.GetHotel)
	apiV1Hotel.Post("/delete", handlers.HotelHandler.DeleteHotel)
	apiV1Hotel.Post("/edit", handlers.HotelHandler.UpdateHotel)

	apiV1HotelPublic := apiV1.Group("/hotel/pub")
	apiV1HotelPublic.Post("/list", handlers.HotelHandler.GetHotelsPublic)
	apiV1HotelPublic.Post("/get", handlers.HotelHandler.GetHotelPublic)

	apiV1Room := apiV1.Group("/room", middlewares.VerifyJwtTokenMiddleware)
	apiV1Room.Post("/add", handlers.RoomHandler.AddRoom)
	apiV1Room.Post("/list", handlers.RoomHandler.GetRooms)
	apiV1Room.Post("/get", handlers.RoomHandler.GetRoom)
	apiV1Room.Post("/delete", handlers.RoomHandler.DeleteRoom)
	apiV1Room.Post("/edit", handlers.RoomHandler.UpdateRoom)

	apiV1RoomPublic := apiV1.Group("/room/pub")
	apiV1RoomPublic.Post("/list", handlers.RoomHandler.GetRoomsPublic)
	apiV1RoomPublic.Post("/get", handlers.RoomHandler.GetRoomPublic)

	app.Listen(*listenAddr)
}
