//THis File is AUTO-GENERATED
package characteristic

const TypeHue = "13"

type Hue struct {
	*Float
}

func NewHue() *Hue {

	char := NewFloat(TypeHue)
	char.Format = FormatFloat
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(0)
	char.SetMaxValue(360)
	char.SetStepValue(1)
	char.Unit = UnitArcDegrees
	char.SetValue(0)
	return &Hue{char}
}
