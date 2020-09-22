//THis File is AUTO-GENERATED
package characteristic

const (
	TargetSlatStateAuto   int = 1
	TargetSlatStateManual int = 0
)
const TypeTargetSlatState = "BE"

type TargetSlatState struct {
	*Int
}

func NewTargetSlatState() *TargetSlatState {

	char := NewInt(TypeTargetSlatState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &TargetSlatState{char}
}
