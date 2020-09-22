//THis File is AUTO-GENERATED
package characteristic

const (
	CurrentMediaStatePause   int = 1
	CurrentMediaStatePlay    int = 0
	CurrentMediaStateStop    int = 2
	CurrentMediaStateUnknown int = 3
)
const TypeCurrentMediaState = "E0"

type CurrentMediaState struct {
	*Int
}

func NewCurrentMediaState() *CurrentMediaState {

	char := NewInt(TypeCurrentMediaState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &CurrentMediaState{char}
}
