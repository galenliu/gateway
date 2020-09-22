//THis File is AUTO-GENERATED
package characteristic

const TypeColorTemperature = "CE"

type ColorTemperature struct {
	*Int
}

func NewColorTemperature() *ColorTemperature {

	char := NewInt(TypeColorTemperature)
	char.Format = FormatUint32
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(140)
	char.SetMaxValue(500)
	char.SetStepValue(1)

	char.SetValue(140)
	return &ColorTemperature{char}
}
