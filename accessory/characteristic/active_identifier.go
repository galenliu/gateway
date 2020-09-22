//THis File is AUTO-GENERATED
package characteristic

const TypeActiveIdentifier = "E7"

type ActiveIdentifier struct {
	*Int
}

func NewActiveIdentifier() *ActiveIdentifier {

	char := NewInt(TypeActiveIdentifier)
	char.Format = FormatUint32
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &ActiveIdentifier{char}
}
