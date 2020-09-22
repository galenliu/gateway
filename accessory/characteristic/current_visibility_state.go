//THis File is AUTO-GENERATED
package characteristic

const (
	CurrentVisibilityStateHidden int = 1
	CurrentVisibilityStateShown  int = 0
	CurrentVisibilityStateState2 int = 2
	CurrentVisibilityStateState3 int = 3
)
const TypeCurrentVisibilityState = "135"

type CurrentVisibilityState struct {
	*Int
}

func NewCurrentVisibilityState() *CurrentVisibilityState {

	char := NewInt(TypeCurrentVisibilityState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &CurrentVisibilityState{char}
}
