//THis File is AUTO-GENERATED
package characteristic

const TypeTargetPosition = "7C"

type TargetPosition struct {
	*Int
}

func NewTargetPosition() *TargetPosition {

	char := NewInt(TypeTargetPosition)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(0)
	char.SetMaxValue(100)
	char.SetStepValue(1)
	char.Unit = UnitPercentage
	char.SetValue(0)
	return &TargetPosition{char}
}
