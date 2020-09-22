//THis File is AUTO-GENERATED
package characteristic

const TypeVolume = "119"

type Volume struct {
	*Int
}

func NewVolume() *Volume {

	char := NewInt(TypeVolume)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(0)
	char.SetMaxValue(100)
	char.SetStepValue(1)
	char.Unit = UnitPercentage
	char.SetValue(0)
	return &Volume{char}
}
