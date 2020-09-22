//THis File is AUTO-GENERATED
package characteristic

const (
	CurrentFanStateBlowingAir int = 2
	CurrentFanStateIdle       int = 1
	CurrentFanStateInactive   int = 0
)
const TypeCurrentFanState = "AF"

type CurrentFanState struct {
	*Int
}

func NewCurrentFanState() *CurrentFanState {

	char := NewInt(TypeCurrentFanState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &CurrentFanState{char}
}
