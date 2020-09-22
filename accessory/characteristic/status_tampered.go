//THis File is AUTO-GENERATED
package characteristic

const (
	StatusTamperedNotTampered int = 0
	StatusTamperedTampered    int = 1
)
const TypeStatusTampered = "7A"

type StatusTampered struct {
	*Int
}

func NewStatusTampered() *StatusTampered {

	char := NewInt(TypeStatusTampered)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &StatusTampered{char}
}
