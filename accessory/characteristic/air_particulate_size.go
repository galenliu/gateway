//THis File is AUTO-GENERATED
package characteristic

const (
	AirParticulateSize10μm  int = 1
	AirParticulateSize2_5μm int = 0
)
const TypeAirParticulateSize = "65"

type AirParticulateSize struct {
	*Int
}

func NewAirParticulateSize() *AirParticulateSize {

	char := NewInt(TypeAirParticulateSize)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &AirParticulateSize{char}
}
