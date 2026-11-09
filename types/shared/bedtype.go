package shared

type BedType string

const (
	BedSingle BedType = "single"
	BedDouble BedType = "double"
	BedQueen  BedType = "queen"
	BedKing   BedType = "king"
	BedTwin   BedType = "twin"
	BedBunk   BedType = "bunk"
)

func (b BedType) Validate() bool {
	switch b {
	case BedSingle, BedDouble, BedQueen, BedKing, BedTwin, BedBunk:
		return true
	}
	return false
}
