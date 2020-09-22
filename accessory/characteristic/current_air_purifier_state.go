//THis File is AUTO-GENERATED
package characteristic

const (
	CurrentAirPurifierStateIdle         int = 1
	CurrentAirPurifierStateInactive     int = 0
	CurrentAirPurifierStatePurifyingAir int = 2
)
const TypeCurrentAirPurifierState = "A9"

type CurrentAirPurifierState struct {
	*Int
}

func NewCurrentAirPurifierState() *CurrentAirPurifierState {

	char := NewInt(TypeCurrentAirPurifierState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &CurrentAirPurifierState{char}
}
