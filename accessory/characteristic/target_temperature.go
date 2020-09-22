//THis File is AUTO-GENERATED
package characteristic

const TypeTargetTemperature = "35"

type TargetTemperature struct {
	*Float
}

func NewTargetTemperature() *TargetTemperature {

	char := NewFloat(TypeTargetTemperature)
	char.Format = FormatFloat
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(10)
	char.SetMaxValue(38)
	char.SetStepValue(0.1)
	char.Unit = UnitCelsius
	char.SetValue(10)
	return &TargetTemperature{char}
}
