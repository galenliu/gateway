//THis File is AUTO-GENERATED
package characteristic

const (
	LockCurrentStateJammed    int = 2
	LockCurrentStateSecured   int = 1
	LockCurrentStateUnknown   int = 3
	LockCurrentStateUnsecured int = 0
)
const TypeLockCurrentState = "1D"

type LockCurrentState struct {
	*Int
}

func NewLockCurrentState() *LockCurrentState {

	char := NewInt(TypeLockCurrentState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &LockCurrentState{char}
}
