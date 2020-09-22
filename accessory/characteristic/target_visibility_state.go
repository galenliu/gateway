//THis File is AUTO-GENERATED
package characteristic

const (
	TargetVisibilityStateHidden int = 1
	TargetVisibilityStateShown  int = 0
)
const TypeTargetVisibilityState = "134"

type TargetVisibilityState struct {
	*Int
}

func NewTargetVisibilityState() *TargetVisibilityState {

	char := NewInt(TypeTargetVisibilityState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &TargetVisibilityState{char}
}
