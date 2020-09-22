//THis File is AUTO-GENERATED
package characteristic

const TypeTargetVerticalTiltAngle = "7D"

type TargetVerticalTiltAngle struct {
	*Int
}

func NewTargetVerticalTiltAngle() *TargetVerticalTiltAngle {

	char := NewInt(TypeTargetVerticalTiltAngle)
	char.Format = FormatInt
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(-90)
	char.SetMaxValue(90)
	char.SetStepValue(1)
	char.Unit = UnitArcDegrees
	char.SetValue(-90)
	return &TargetVerticalTiltAngle{char}
}
