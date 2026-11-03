package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// todo!! => must move to dotenv
const dbName = "hotelReservation"

var (
	ctx        = context.Background()
	client     = db.NewMongoDbClient()
	hotelStore = db.NewMongoHotelStore(client, dbName)
	roomStore  = db.NewMongoRoomStore(client, hotelStore, dbName)
)

type seedHotelAndRoomParams struct {
	name     string
	location string
	rate     int
}

func seedHotelAndRooms(params *seedHotelAndRoomParams) {
	hotel := types.NewHotel(&types.CreateHotelParams{Name: params.name, Location: params.location , Rate: params.rate})

	rooms := []*types.Room{
		types.NewRoom(&types.CreateRoomParams{
			Type:      types.SingleRoomType,
			BasePrice: 99.9,
			Price:     0.0,
			HotelId:   primitive.NilObjectID,
		}),
		types.NewRoom(&types.CreateRoomParams{
			Type:      types.DeluxeRoomType,
			BasePrice: 199.9,
			Price:     0.0,
			HotelId:   primitive.NilObjectID,
		}),
		types.NewRoom(&types.CreateRoomParams{
			Type:      types.SeaSideRoomType,
			BasePrice: 122.9,
			Price:     0.0,
			HotelId:   primitive.NilObjectID,
		}),
	}

	hotelInsertResult, err := hotelStore.InsertHotel(ctx, hotel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\n", params.name, params.location)
	fmt.Printf("%+v\n", hotelInsertResult)
	for _, room := range rooms {
		room.HotelId = hotelInsertResult.Id
		roomInsertResult, err := roomStore.InsertRoom(ctx, room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", roomInsertResult)
	}
}

func main() {
	hotelStore.Drop(ctx)
	roomStore.Drop(ctx)

	seedHotelAndRooms(&seedHotelAndRoomParams{name: "Espinas Palac Hotel", location: "Tehran" , rate: 3})
	seedHotelAndRooms(&seedHotelAndRoomParams{name: "Abbasi Hotel", location: "Isfahan" , rate: 4})
	seedHotelAndRooms(&seedHotelAndRoomParams{name: "Toranj Marine Hotel", location: "Kish Island" , rate: 1})
}
