//THis File is AUTO-GENERATED
package characteristic

const (
	ActiveActive   int = 1
	ActiveInactive int = 0
)
const TypeActive = "B0"

type Active struct {
	*Int
}

func NewActive() *Active {

	char := NewInt(TypeActive)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &Active{char}
}
