//THis File is AUTO-GENERATED
package characteristic

const TypeCurrentRelativeHumidity = "10"

type CurrentRelativeHumidity struct {
	*Float
}

func NewCurrentRelativeHumidity() *CurrentRelativeHumidity {

	char := NewFloat(TypeCurrentRelativeHumidity)
	char.Format = FormatFloat
	char.Perms = []string{PermRead, PermEvents}
	char.SetMinValue(0)
	char.SetMaxValue(100)
	char.SetStepValue(1)
	char.Unit = UnitPercentage
	char.SetValue(0)
	return &CurrentRelativeHumidity{char}
}
