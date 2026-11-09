package shared

type Permission string

const (
	PermHotelCreate Permission = "hotel:create"
	PermHotelUpdate Permission = "hotel:update"
	PermHotelDelete Permission = "hotel:delete"

	PermRoomCreate Permission = "room:create"
	PermRoomUpdate Permission = "room:update"
	PermRoomDelete Permission = "room:delete"

	PermBookingCreate Permission = "booking:create"
	PermBookingRead   Permission = "booking:read"
	PermBookingCancel Permission = "booking:cancel"
	PermBookingManage Permission = "booking:manage"

	PermUserManage Permission = "user:manage"
	PermRoleManage Permission = "role:manage"
)

func (p Permission) Validate() bool {
	switch p {
	case PermHotelCreate, PermHotelUpdate, PermHotelDelete,
		PermRoomCreate, PermRoomUpdate, PermRoomDelete,
		PermBookingCreate, PermBookingRead, PermBookingCancel, PermBookingManage,
		PermUserManage, PermRoleManage:
		return true
	}
	return false
}
