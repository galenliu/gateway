//THis File is AUTO-GENERATED
package characteristic

const (
	LockLastKnownActionSecuredPhysically_Exterior   int = 2
	LockLastKnownActionSecuredPhysically_Interior   int = 0
	LockLastKnownActionSecuredRemotely              int = 6
	LockLastKnownActionSecuredbyAutoSecureTimeout   int = 8
	LockLastKnownActionSecuredbyKeypad              int = 4
	LockLastKnownActionUnsecuredPhysically_Exterior int = 3
	LockLastKnownActionUnsecuredPhysically_Interior int = 1
	LockLastKnownActionUnsecuredRemotely            int = 7
	LockLastKnownActionUnsecuredbyKeypad            int = 5
)
const TypeLockLastKnownAction = "1C"

type LockLastKnownAction struct {
	*Int
}

func NewLockLastKnownAction() *LockLastKnownAction {

	char := NewInt(TypeLockLastKnownAction)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &LockLastKnownAction{char}
}
