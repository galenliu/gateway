//THis File is AUTO-GENERATED
package characteristic

const TypeRotationSpeed = "29"

type RotationSpeed struct {
	*Float
}

func NewRotationSpeed() *RotationSpeed {

	char := NewFloat(TypeRotationSpeed)
	char.Format = FormatFloat
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(0)
	char.SetMaxValue(100)
	char.SetStepValue(1)
	char.Unit = UnitPercentage
	char.SetValue(0)
	return &RotationSpeed{char}
}
