//THis File is AUTO-GENERATED
package characteristic

const (
	InputDeviceTypeAudioSystem int = 5
	InputDeviceTypeOther       int = 0
	InputDeviceTypePlayback    int = 4
	InputDeviceTypeRecording   int = 2
	InputDeviceTypeTV          int = 1
	InputDeviceTypeTuner       int = 3
)
const TypeInputDeviceType = "DC"

type InputDeviceType struct {
	*Int
}

func NewInputDeviceType() *InputDeviceType {

	char := NewInt(TypeInputDeviceType)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &InputDeviceType{char}
}
