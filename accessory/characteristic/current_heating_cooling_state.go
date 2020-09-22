//THis File is AUTO-GENERATED
package characteristic

const (
	CurrentHeatingCoolingStateCool int = 2
	CurrentHeatingCoolingStateHeat int = 1
	CurrentHeatingCoolingStateOff  int = 0
)
const TypeCurrentHeatingCoolingState = "F"

type CurrentHeatingCoolingState struct {
	*Int
}

func NewCurrentHeatingCoolingState() *CurrentHeatingCoolingState {

	char := NewInt(TypeCurrentHeatingCoolingState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &CurrentHeatingCoolingState{char}
}
