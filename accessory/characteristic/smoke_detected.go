//THis File is AUTO-GENERATED
package characteristic

const (
	SmokeDetectedSmokeDetected    int = 1
	SmokeDetectedSmokeNotDetected int = 0
)
const TypeSmokeDetected = "76"

type SmokeDetected struct {
	*Int
}

func NewSmokeDetected() *SmokeDetected {

	char := NewInt(TypeSmokeDetected)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &SmokeDetected{char}
}
