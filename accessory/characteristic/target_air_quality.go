//THis File is AUTO-GENERATED
package characteristic

const (
	TargetAirQualityExcellent int = 0
	TargetAirQualityFair      int = 2
	TargetAirQualityGood      int = 1
)
const TypeTargetAirQuality = "AE"

type TargetAirQuality struct {
	*Int
}

func NewTargetAirQuality() *TargetAirQuality {

	char := NewInt(TypeTargetAirQuality)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &TargetAirQuality{char}
}
