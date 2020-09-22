//THis File is AUTO-GENERATED
package characteristic

const (
	StatusFaultGeneralFault int = 1
	StatusFaultNoFault      int = 0
)
const TypeStatusFault = "77"

type StatusFault struct {
	*Int
}

func NewStatusFault() *StatusFault {

	char := NewInt(TypeStatusFault)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &StatusFault{char}
}
