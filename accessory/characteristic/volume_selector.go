//THis File is AUTO-GENERATED
package characteristic

const (
	VolumeSelectorDecrement int = 1
	VolumeSelectorIncrement int = 0
)
const TypeVolumeSelector = "EA"

type VolumeSelector struct {
	*Int
}

func NewVolumeSelector() *VolumeSelector {

	char := NewInt(TypeVolumeSelector)
	char.Format = FormatUint8
	char.Perms = []string{PermWrite}

	char.SetValue(0)
	return &VolumeSelector{char}
}
