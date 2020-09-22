//THis File is AUTO-GENERATED
package characteristic

const (
	CarbonDioxideDetectedCO2LevelsAbnormal int = 1
	CarbonDioxideDetectedCO2LevelsNormal   int = 0
)
const TypeCarbonDioxideDetected = "92"

type CarbonDioxideDetected struct {
	*Int
}

func NewCarbonDioxideDetected() *CarbonDioxideDetected {

	char := NewInt(TypeCarbonDioxideDetected)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &CarbonDioxideDetected{char}
}
