//THis File is AUTO-GENERATED
package characteristic

const TypeCurrentAmbientLightLevel = "6B"

type CurrentAmbientLightLevel struct {
	*Float
}

func NewCurrentAmbientLightLevel() *CurrentAmbientLightLevel {

	char := NewFloat(TypeCurrentAmbientLightLevel)
	char.Format = FormatFloat
	char.Perms = []string{PermRead, PermEvents}
	char.SetMinValue(0.0001)
	char.SetMaxValue(100000)

	char.Unit = UnitLux
	char.SetValue(0.0001)
	return &CurrentAmbientLightLevel{char}
}
