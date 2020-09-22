//THis File is AUTO-GENERATED
package characteristic

const (
	RotationDirectionClockwise         int = 0
	RotationDirectionCounter_clockwise int = 1
)
const TypeRotationDirection = "28"

type RotationDirection struct {
	*Int
}

func NewRotationDirection() *RotationDirection {

	char := NewInt(TypeRotationDirection)
	char.Format = FormatInt
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &RotationDirection{char}
}
