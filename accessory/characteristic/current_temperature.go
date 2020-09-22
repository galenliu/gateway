//THis File is AUTO-GENERATED
package characteristic

const TypeCurrentTemperature = "11"

type CurrentTemperature struct {
	*Float
}

func NewCurrentTemperature() *CurrentTemperature {

	char := NewFloat(TypeCurrentTemperature)
	char.Format = FormatFloat
	char.Perms = []string{PermRead, PermEvents}
	char.SetMinValue(-273.1)
	char.SetMaxValue(1000)
	char.SetStepValue(0.1)
	char.Unit = UnitCelsius
	char.SetValue(-273.1)
	return &CurrentTemperature{char}
}
