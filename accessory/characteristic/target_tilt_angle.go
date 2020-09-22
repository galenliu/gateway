//THis File is AUTO-GENERATED
package characteristic

const TypeTargetTiltAngle = "C2"

type TargetTiltAngle struct {
	*Int
}

func NewTargetTiltAngle() *TargetTiltAngle {

	char := NewInt(TypeTargetTiltAngle)
	char.Format = FormatInt
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(-90)
	char.SetMaxValue(90)
	char.SetStepValue(1)
	char.Unit = UnitArcDegrees
	char.SetValue(-90)
	return &TargetTiltAngle{char}
}
