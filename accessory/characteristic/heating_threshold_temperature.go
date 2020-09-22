//THis File is AUTO-GENERATED
package characteristic

const TypeHeatingThresholdTemperature = "12"

type HeatingThresholdTemperature struct {
	*Float
}

func NewHeatingThresholdTemperature() *HeatingThresholdTemperature {

	char := NewFloat(TypeHeatingThresholdTemperature)
	char.Format = FormatFloat
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(0)
	char.SetMaxValue(25)
	char.SetStepValue(0.1)
	char.Unit = UnitCelsius
	char.SetValue(0)
	return &HeatingThresholdTemperature{char}
}
