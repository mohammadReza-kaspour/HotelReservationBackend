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

func main() {
	ctx := context.Background()
	client := db.NewMongoDbClient()

	hotelStore := db.NewMongoHotelStore(client, dbName)
	roomStore := db.NewMongoRoomStore(client, hotelStore, dbName)

	hotel := types.NewHotel(&types.CreateHotelParams{Name: "Espinas Palace", Location: "Tehran"})

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

	fmt.Println(hotelInsertResult)
	for _, room := range rooms {
		room.HotelId = hotelInsertResult.Id
		roomInsertResult, err := roomStore.InsertRoom(ctx, room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(roomInsertResult)
	}
}
