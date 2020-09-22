//THis File is AUTO-GENERATED
package characteristic

const TypeCurrentPosition = "6D"

type CurrentPosition struct {
	*Int
}

func NewCurrentPosition() *CurrentPosition {

	char := NewInt(TypeCurrentPosition)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}
	char.SetMinValue(0)
	char.SetMaxValue(100)
	char.SetStepValue(1)
	char.Unit = UnitPercentage
	char.SetValue(0)
	return &CurrentPosition{char}
}
