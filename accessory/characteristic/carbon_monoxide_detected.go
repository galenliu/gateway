//THis File is AUTO-GENERATED
package characteristic

const (
	CarbonMonoxideDetectedCOLevelsAbnormal int = 1
	CarbonMonoxideDetectedCOLevelsNormal   int = 0
)
const TypeCarbonMonoxideDetected = "69"

type CarbonMonoxideDetected struct {
	*Int
}

func NewCarbonMonoxideDetected() *CarbonMonoxideDetected {

	char := NewInt(TypeCarbonMonoxideDetected)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &CarbonMonoxideDetected{char}
}
