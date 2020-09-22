//THis File is AUTO-GENERATED
package characteristic

const (
	AirQualityExcellent int = 1
	AirQualityFair      int = 3
	AirQualityGood      int = 2
	AirQualityInferior  int = 4
	AirQualityPoor      int = 5
	AirQualityUnknown   int = 0
)
const TypeAirQuality = "95"

type AirQuality struct {
	*Int
}

func NewAirQuality() *AirQuality {

	char := NewInt(TypeAirQuality)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &AirQuality{char}
}
