package shared

type Amenity string

const (
	// Hotel level
	AmenityWifi        Amenity = "wifi"
	AmenityParking     Amenity = "parking"
	AmenityPool        Amenity = "pool"
	AmenityGym         Amenity = "gym"
	AmenitySpa         Amenity = "spa"
	AmenityRestaurant  Amenity = "restaurant"
	AmenityBar         Amenity = "bar"
	AmenityPetFriendly Amenity = "pet_friendly"
	AmenityAirportPick Amenity = "airport_pickup"

	// Room level
	AmenityAC        Amenity = "ac"
	AmenityTV        Amenity = "tv"
	AmenityMinibar   Amenity = "minibar"
	AmenityBalcony   Amenity = "balcony"
	AmenityBathtub   Amenity = "bathtub"
	AmenityKitchen   Amenity = "kitchen"
	AmenityDesk      Amenity = "desk"
	AmenitySafe      Amenity = "safe"
	AmenityBreakfast Amenity = "breakfast"
)

func (a Amenity) Validate() bool {
	switch a {
	case AmenityWifi, AmenityParking, AmenityPool, AmenityGym, AmenitySpa,
		AmenityRestaurant, AmenityBar, AmenityPetFriendly, AmenityAirportPick,
		AmenityAC, AmenityTV, AmenityMinibar, AmenityBalcony, AmenityBathtub,
		AmenityKitchen, AmenityDesk, AmenitySafe, AmenityBreakfast:
		return true
	}
	return false
}
