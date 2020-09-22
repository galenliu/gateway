//THis File is AUTO-GENERATED
package characteristic

const (
	TargetHeaterCoolerStateAuto int = 0
	TargetHeaterCoolerStateCool int = 2
	TargetHeaterCoolerStateHeat int = 1
)
const TypeTargetHeaterCoolerState = "B2"

type TargetHeaterCoolerState struct {
	*Int
}

func NewTargetHeaterCoolerState() *TargetHeaterCoolerState {

	char := NewInt(TypeTargetHeaterCoolerState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &TargetHeaterCoolerState{char}
}
