//THis File is AUTO-GENERATED
package characteristic

const TypeLockManagementAutoSecurityTimeout = "1A"

type LockManagementAutoSecurityTimeout struct {
	*Int
}

func NewLockManagementAutoSecurityTimeout() *LockManagementAutoSecurityTimeout {

	char := NewInt(TypeLockManagementAutoSecurityTimeout)
	char.Format = FormatUint32
	char.Perms = []string{PermRead, PermWrite, PermEvents}

	char.Unit = UnitSeconds
	char.SetValue(0)
	return &LockManagementAutoSecurityTimeout{char}
}
