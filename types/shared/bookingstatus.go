package shared

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingCancelled BookingStatus = "cancelled"
	BookingCompleted BookingStatus = "completed"
)

func (s BookingStatus) Validate() bool {
	switch s {
	case BookingPending, BookingConfirmed, BookingCancelled, BookingCompleted:
		return true
	}
	return false
}
