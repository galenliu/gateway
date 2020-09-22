//THis File is AUTO-GENERATED
package characteristic

const (
	IsConfiguredConfigured    int = 1
	IsConfiguredNotConfigured int = 0
)
const TypeIsConfigured = "D6"

type IsConfigured struct {
	*Int
}

func NewIsConfigured() *IsConfigured {

	char := NewInt(TypeIsConfigured)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &IsConfigured{char}
}
