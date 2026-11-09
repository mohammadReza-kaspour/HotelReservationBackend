package shared

type RoomView string

const (
	ViewNone     RoomView = "none"
	ViewSea      RoomView = "sea"
	ViewCity     RoomView = "city"
	ViewGarden   RoomView = "garden"
	ViewPool     RoomView = "pool"
	ViewMountain RoomView = "mountain"
)

func (v RoomView) Validate() bool {
	switch v {
	case ViewNone, ViewSea, ViewCity, ViewGarden, ViewPool, ViewMountain:
		return true
	}
	return false
}
