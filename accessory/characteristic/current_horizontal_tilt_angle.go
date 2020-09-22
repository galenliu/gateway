//THis File is AUTO-GENERATED
package characteristic

const TypeCurrentHorizontalTiltAngle = "6C"

type CurrentHorizontalTiltAngle struct {
	*Int
}

func NewCurrentHorizontalTiltAngle() *CurrentHorizontalTiltAngle {

	char := NewInt(TypeCurrentHorizontalTiltAngle)
	char.Format = FormatInt
	char.Perms = []string{PermRead, PermEvents}
	char.SetMinValue(-90)
	char.SetMaxValue(90)
	char.SetStepValue(1)
	char.Unit = UnitArcDegrees
	char.SetValue(-90)
	return &CurrentHorizontalTiltAngle{char}
}
