//THis File is AUTO-GENERATED
package characteristic

const (
	ValveTypeGenericvalve int = 0
	ValveTypeIrrigation   int = 1
	ValveTypeShowerhead   int = 2
	ValveTypeWaterfaucet  int = 3
)
const TypeValveType = "D5"

type ValveType struct {
	*Int
}

func NewValveType() *ValveType {

	char := NewInt(TypeValveType)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &ValveType{char}
}
