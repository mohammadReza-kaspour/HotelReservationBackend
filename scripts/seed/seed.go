package main

import (
	"context"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/db"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/config"
)

var (
	ctx        = context.Background()
	client     = db.NewMongoDbClient()
	hotelStore = db.NewMongoHotelStore(client, config.DbName)
	roomStore  = db.NewMongoRoomStore(client, hotelStore, config.DbName)
)

type seedHotelAndRoomParams struct {
	name     string
	location string
	rate     int
}

func main() {
}
