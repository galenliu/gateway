//THis File is AUTO-GENERATED
package characteristic

const (
	TargetHeatingCoolingStateAuto int = 3
	TargetHeatingCoolingStateCool int = 2
	TargetHeatingCoolingStateHeat int = 1
	TargetHeatingCoolingStateOff  int = 0
)
const TypeTargetHeatingCoolingState = "33"

type TargetHeatingCoolingState struct {
	*Int
}

func NewTargetHeatingCoolingState() *TargetHeatingCoolingState {

	char := NewInt(TypeTargetHeatingCoolingState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &TargetHeatingCoolingState{char}
}
