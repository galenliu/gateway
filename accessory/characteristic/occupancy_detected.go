//THis File is AUTO-GENERATED
package characteristic

const (
	OccupancyDetectedOccupancyDetected    int = 1
	OccupancyDetectedOccupancyNotDetected int = 0
)
const TypeOccupancyDetected = "71"

type OccupancyDetected struct {
	*Int
}

func NewOccupancyDetected() *OccupancyDetected {

	char := NewInt(TypeOccupancyDetected)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &OccupancyDetected{char}
}
