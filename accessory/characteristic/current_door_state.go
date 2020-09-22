//THis File is AUTO-GENERATED
package characteristic

const (
	CurrentDoorStateClosed  int = 1
	CurrentDoorStateClosing int = 3
	CurrentDoorStateOpen    int = 0
	CurrentDoorStateOpening int = 2
	CurrentDoorStateStopped int = 4
)
const TypeCurrentDoorState = "E"

type CurrentDoorState struct {
	*Int
}

func NewCurrentDoorState() *CurrentDoorState {

	char := NewInt(TypeCurrentDoorState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &CurrentDoorState{char}
}
