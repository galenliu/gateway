//THis File is AUTO-GENERATED
package characteristic

const (
	InUseInuse    int = 1
	InUseNotinuse int = 0
)
const TypeInUse = "D2"

type InUse struct {
	*Int
}

func NewInUse() *InUse {

	char := NewInt(TypeInUse)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &InUse{char}
}
