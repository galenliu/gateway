//THis File is AUTO-GENERATED
package characteristic

const (
	SecuritySystemCurrentStateAlarmTriggered int = 4
	SecuritySystemCurrentStateAwayArm        int = 1
	SecuritySystemCurrentStateDisarmed       int = 3
	SecuritySystemCurrentStateNightArm       int = 2
	SecuritySystemCurrentStateStayArm        int = 0
)
const TypeSecuritySystemCurrentState = "66"

type SecuritySystemCurrentState struct {
	*Int
}

func NewSecuritySystemCurrentState() *SecuritySystemCurrentState {

	char := NewInt(TypeSecuritySystemCurrentState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &SecuritySystemCurrentState{char}
}
