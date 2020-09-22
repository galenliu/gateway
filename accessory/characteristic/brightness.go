//THis File is AUTO-GENERATED
package characteristic

const TypeBrightness = "8"

type Brightness struct {
	*Int
}

func NewBrightness() *Brightness {

	char := NewInt(TypeBrightness)
	char.Format = FormatInt
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(0)
	char.SetMaxValue(100)
	char.SetStepValue(1)
	char.Unit = UnitPercentage
	char.SetValue(0)
	return &Brightness{char}
}
