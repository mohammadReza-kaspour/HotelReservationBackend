package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/api"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/middlewares/auth"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
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
	stringifyInput , err := json.Marshal(map[string]any{"firstName":"Jane","lastName":"Doe","email":"jane.doe@example.com","password":"P@ssw0rd123"})
	if err != nil {log.Fatal(err)}

	var parsedInput *types.CreateUserParams;
	if err = json.NewDecoder(bytes.NewReader(stringifyInput)).Decode(parsedInput);err != nil {
		log.Fatal(err)
	}
	fmt.Println(parsedInput)
	log.Fatal("finished")


	listenAddr := flag.String("listenAddr", ":5000", "the listen addr of api server")
	flag.Parse()

	//create mongo client
	client := db.NewMongoDbClient()

	//user store and handlers initialization
	userStore := db.NewMongoUserStore(client, dbName)
	userHandler := api.NewUserHandler(userStore)
	authHandler := api.NewAuthHandler(userStore)

	hotelStore := db.NewMongoHotelStore(client, dbName)
	hotelHandler := api.NewHotelHandler(hotelStore)

	app := fiber.New(config)

	app.Post("/signup", authHandler.SignUp)
	app.Post("/signin", authHandler.SignIn)              //not implemented yet
	app.Post("/refresh-token", authHandler.RefreshToken) //not implemented yet
	// app.Get("/developer", func(c *fiber.Ctx) error {
	// 	return errors.New("error occured")
	// 	return c.Status(200).SendString("ok")
	// })

	apiv1 := app.Group("/api/v1", auth.VerifyJwtTokenMiddleware)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user", userHandler.HandleGetUser)
	// apiv1.Post("/user", userHandler.HandlePostUser) // !!Obsolete after adding jwt token
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)

	app.Listen(*listenAddr)
}
