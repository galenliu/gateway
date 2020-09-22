//THis File is AUTO-GENERATED
package characteristic

const (
	StatusLowBatteryBatteryLevelLow    int = 1
	StatusLowBatteryBatteryLevelNormal int = 0
)
const TypeStatusLowBattery = "79"

type StatusLowBattery struct {
	*Int
}

func NewStatusLowBattery() *StatusLowBattery {

	char := NewInt(TypeStatusLowBattery)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &StatusLowBattery{char}
}
