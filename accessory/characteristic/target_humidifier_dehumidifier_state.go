//THis File is AUTO-GENERATED
package characteristic

const (
	TargetHumidifierDehumidifierStateDehumidifier             int = 2
	TargetHumidifierDehumidifierStateHumidifier               int = 1
	TargetHumidifierDehumidifierStateHumidifierorDehumidifier int = 0
)
const TypeTargetHumidifierDehumidifierState = "B4"

type TargetHumidifierDehumidifierState struct {
	*Int
}

func NewTargetHumidifierDehumidifierState() *TargetHumidifierDehumidifierState {

	char := NewInt(TypeTargetHumidifierDehumidifierState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &TargetHumidifierDehumidifierState{char}
}
