//THis File is AUTO-GENERATED
package characteristic

const (
	TargetAirPurifierStateAuto   int = 1
	TargetAirPurifierStateManual int = 0
)
const TypeTargetAirPurifierState = "A8"

type TargetAirPurifierState struct {
	*Int
}

func NewTargetAirPurifierState() *TargetAirPurifierState {

	char := NewInt(TypeTargetAirPurifierState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &TargetAirPurifierState{char}
}
