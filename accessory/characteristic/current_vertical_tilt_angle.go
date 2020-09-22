//THis File is AUTO-GENERATED
package characteristic

const TypeCurrentVerticalTiltAngle = "6E"

type CurrentVerticalTiltAngle struct {
	*Int
}

func NewCurrentVerticalTiltAngle() *CurrentVerticalTiltAngle {

	char := NewInt(TypeCurrentVerticalTiltAngle)
	char.Format = FormatInt
	char.Perms = []string{PermRead, PermEvents}
	char.SetMinValue(-90)
	char.SetMaxValue(90)
	char.SetStepValue(1)
	char.Unit = UnitArcDegrees
	char.SetValue(-90)
	return &CurrentVerticalTiltAngle{char}
}
