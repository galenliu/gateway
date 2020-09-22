//THis File is AUTO-GENERATED
package characteristic

const (
	CurrentHumidifierDehumidifierStateDehumidifying int = 3
	CurrentHumidifierDehumidifierStateHumidifying   int = 2
	CurrentHumidifierDehumidifierStateIdle          int = 1
	CurrentHumidifierDehumidifierStateInactive      int = 0
)
const TypeCurrentHumidifierDehumidifierState = "B3"

type CurrentHumidifierDehumidifierState struct {
	*Int
}

func NewCurrentHumidifierDehumidifierState() *CurrentHumidifierDehumidifierState {

	char := NewInt(TypeCurrentHumidifierDehumidifierState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &CurrentHumidifierDehumidifierState{char}
}
