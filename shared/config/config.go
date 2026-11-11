package config

import "time"

const (
	//DB
	DbConnectionTimout      = time.Second * 10
	DbName                  = "hotelReservation"
	DbName_Test             = "hotelReservation_test"
	DbUri                   = "mongodb://localhost:27017"
	DbUserCollectionName    = "users"
	DbHotelCollectionName   = "hotels"
	DbRoomCollectionName    = "rooms"
	DbSessionCollectionName = "sessions"
	DbRoleCollectionName    = "roles"
	DbBookingCollectionName = "bookings"

	//AUTH
	AccessTokenExpireDelay  = time.Hour * 24 * 365
	RefreshTokenExpireDelay = time.Hour * 24 * 2
	SessionExpireTime       = time.Hour * 24 * 30
)
