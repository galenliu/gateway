//THis File is AUTO-GENERATED
package characteristic

const (
	FilterChangeIndicationChangeFilter int = 1
	FilterChangeIndicationFilterOK     int = 0
)
const TypeFilterChangeIndication = "AC"

type FilterChangeIndication struct {
	*Int
}

func NewFilterChangeIndication() *FilterChangeIndication {

	char := NewInt(TypeFilterChangeIndication)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &FilterChangeIndication{char}
}
