package shared

type PaymentStatus string

const (
	PaymentUnpaid   PaymentStatus = "unpaid"
	PaymentPaid     PaymentStatus = "paid"
	PaymentRefunded PaymentStatus = "refunded"
)

func (s PaymentStatus) Validate() bool {
	switch s {
	case PaymentUnpaid, PaymentPaid, PaymentRefunded:
		return true
	}
	return false
}
