//THis File is AUTO-GENERATED
package characteristic

const TypeCurrentTiltAngle = "C1"

type CurrentTiltAngle struct {
	*Int
}

func NewCurrentTiltAngle() *CurrentTiltAngle {

	char := NewInt(TypeCurrentTiltAngle)
	char.Format = FormatInt
	char.Perms = []string{PermRead, PermEvents}
	char.SetMinValue(-90)
	char.SetMaxValue(90)
	char.SetStepValue(1)
	char.Unit = UnitArcDegrees
	char.SetValue(-90)
	return &CurrentTiltAngle{char}
}
