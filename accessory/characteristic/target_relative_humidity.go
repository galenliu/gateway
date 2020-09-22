//THis File is AUTO-GENERATED
package characteristic

const TypeTargetRelativeHumidity = "34"

type TargetRelativeHumidity struct {
	*Float
}

func NewTargetRelativeHumidity() *TargetRelativeHumidity {

	char := NewFloat(TypeTargetRelativeHumidity)
	char.Format = FormatFloat
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(0)
	char.SetMaxValue(100)
	char.SetStepValue(1)
	char.Unit = UnitPercentage
	char.SetValue(0)
	return &TargetRelativeHumidity{char}
}
