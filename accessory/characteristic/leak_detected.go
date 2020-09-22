//THis File is AUTO-GENERATED
package characteristic

const (
	LeakDetectedLeakDetected    int = 1
	LeakDetectedLeakNotDetected int = 0
)
const TypeLeakDetected = "70"

type LeakDetected struct {
	*Int
}

func NewLeakDetected() *LeakDetected {

	char := NewInt(TypeLeakDetected)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &LeakDetected{char}
}
