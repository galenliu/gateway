//THis File is AUTO-GENERATED
package characteristic

const (
	VolumeControlTypeAbsolute            int = 3
	VolumeControlTypeNone                int = 0
	VolumeControlTypeRelative            int = 1
	VolumeControlTypeRelativeWithCurrent int = 2
)
const TypeVolumeControlType = "E9"

type VolumeControlType struct {
	*Int
}

func NewVolumeControlType() *VolumeControlType {

	char := NewInt(TypeVolumeControlType)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &VolumeControlType{char}
}
