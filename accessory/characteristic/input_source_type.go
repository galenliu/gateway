//THis File is AUTO-GENERATED
package characteristic

const (
	InputSourceTypeAirPlay        int = 8
	InputSourceTypeApplication    int = 10
	InputSourceTypeComponentVideo int = 6
	InputSourceTypeCompositeVideo int = 4
	InputSourceTypeDVI            int = 7
	InputSourceTypeHDMI           int = 3
	InputSourceTypeHomeScreen     int = 1
	InputSourceTypeOther          int = 0
	InputSourceTypeSVideo         int = 5
	InputSourceTypeTuner          int = 2
	InputSourceTypeUSB            int = 9
)
const TypeInputSourceType = "DB"

type InputSourceType struct {
	*Int
}

func NewInputSourceType() *InputSourceType {

	char := NewInt(TypeInputSourceType)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &InputSourceType{char}
}
