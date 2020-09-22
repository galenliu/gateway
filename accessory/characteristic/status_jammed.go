//THis File is AUTO-GENERATED
package characteristic

const (
	StatusJammedJammed    int = 1
	StatusJammedNotJammed int = 0
)
const TypeStatusJammed = "78"

type StatusJammed struct {
	*Int
}

func NewStatusJammed() *StatusJammed {

	char := NewInt(TypeStatusJammed)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &StatusJammed{char}
}
