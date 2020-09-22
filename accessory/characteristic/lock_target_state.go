//THis File is AUTO-GENERATED
package characteristic

const (
	LockTargetStateSecured   int = 1
	LockTargetStateUnsecured int = 0
)
const TypeLockTargetState = "1E"

type LockTargetState struct {
	*Int
}

func NewLockTargetState() *LockTargetState {

	char := NewInt(TypeLockTargetState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.SetValue(0)
	return &LockTargetState{char}
}
