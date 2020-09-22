//THis File is AUTO-GENERATED
package characteristic

const TypeSaturation = "2F"

type Saturation struct {
	*Float
}

func NewSaturation() *Saturation {

	char := NewFloat(TypeSaturation)
	char.Format = FormatFloat
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(0)
	char.SetMaxValue(100)
	char.SetStepValue(1)
	char.Unit = UnitPercentage
	char.SetValue(0)
	return &Saturation{char}
}
