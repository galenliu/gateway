//THis File is AUTO-GENERATED
package characteristic

const (
	TargetDoorStateClosed int = 1
	TargetDoorStateOpen   int = 0
)
const TypeTargetDoorState = "32"

type TargetDoorState struct {
	*Int
}

func NewTargetDoorState() *TargetDoorState {

	char := NewInt(TypeTargetDoorState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &TargetDoorState{char}
}
