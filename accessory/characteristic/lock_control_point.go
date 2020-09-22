//THis File is AUTO-GENERATED
package characteristic

const TypeLockControlPoint = "19"

type LockControlPoint struct {
	*Bytes
}

func NewLockControlPoint() *LockControlPoint {

	char := NewBytes(TypeLockControlPoint)
	char.Format = FormatTLV8
	char.Perms = []string{PermWrite}

	char.SetValue([]byte{})
	return &LockControlPoint{char}
}
