//THis File is AUTO-GENERATED
package characteristic

const (
	ClosedCaptionsDisabled int = 0
	ClosedCaptionsEnabled  int = 1
)
const TypeClosedCaptions = "DD"

type ClosedCaptions struct {
	*Int
}

func NewClosedCaptions() *ClosedCaptions {

	char := NewInt(TypeClosedCaptions)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &ClosedCaptions{char}
}
