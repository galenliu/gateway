//THis File is AUTO-GENERATED
package characteristic

const (
	PowerModeSelectionHide int = 1
	PowerModeSelectionShow int = 0
)
const TypePowerModeSelection = "DF"

type PowerModeSelection struct {
	*Int
}

func NewPowerModeSelection() *PowerModeSelection {

	char := NewInt(TypePowerModeSelection)
	char.Format = FormatUint8
	char.Perms = []string{PermWrite}

	char.SetValue(0)
	return &PowerModeSelection{char}
}
