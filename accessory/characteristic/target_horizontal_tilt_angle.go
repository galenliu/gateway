//THis File is AUTO-GENERATED
package characteristic

const TypeTargetHorizontalTiltAngle = "7B"

type TargetHorizontalTiltAngle struct {
	*Int
}

func NewTargetHorizontalTiltAngle() *TargetHorizontalTiltAngle {

	char := NewInt(TypeTargetHorizontalTiltAngle)
	char.Format = FormatInt
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(-90)
	char.SetMaxValue(90)
	char.SetStepValue(1)
	char.Unit = UnitArcDegrees
	char.SetValue(-90)
	return &TargetHorizontalTiltAngle{char}
}
