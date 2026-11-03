package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/api"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
)
//todo!! => must move to dotenv
const dbName = "hotelReservation"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.Status(500).JSON(map[string]string{
			"error": err.Error(),
		})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "the listen addr of api server")
	flag.Parse()

	//create mongo client
	client:= db.NewMongoDbClient()

	//user store and handlers initialization
	userStore := db.NewMongoUserStore(client , dbName)
	userHandler := api.NewUserHandler(userStore)

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	app.Listen(*listenAddr)
}
