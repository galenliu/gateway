//THis File is AUTO-GENERATED
package characteristic

const (
	TargetMediaStatePause int = 1
	TargetMediaStatePlay  int = 0
	TargetMediaStateStop  int = 2
)
const TypeTargetMediaState = "137"

type TargetMediaState struct {
	*Int
}

func NewTargetMediaState() *TargetMediaState {

	char := NewInt(TypeTargetMediaState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &TargetMediaState{char}
}
