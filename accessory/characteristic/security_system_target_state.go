//THis File is AUTO-GENERATED
package characteristic

const (
	SecuritySystemTargetStateAwayArm  int = 1
	SecuritySystemTargetStateDisarm   int = 3
	SecuritySystemTargetStateNightArm int = 2
	SecuritySystemTargetStateStayArm  int = 0
)
const TypeSecuritySystemTargetState = "67"

type SecuritySystemTargetState struct {
	*Int
}

func NewSecuritySystemTargetState() *SecuritySystemTargetState {

	char := NewInt(TypeSecuritySystemTargetState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &SecuritySystemTargetState{char}
}
