//THis File is AUTO-GENERATED
package characteristic

const (
	TargetFanStateAuto   int = 1
	TargetFanStateManual int = 0
)
const TypeTargetFanState = "BF"

type TargetFanState struct {
	*Int
}

func NewTargetFanState() *TargetFanState {

	char := NewInt(TypeTargetFanState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &TargetFanState{char}
}
