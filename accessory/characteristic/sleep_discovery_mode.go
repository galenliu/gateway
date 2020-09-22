//THis File is AUTO-GENERATED
package characteristic

const (
	SleepDiscoveryModeAlwaysDiscoverable int = 1
	SleepDiscoveryModeNotDiscoverable    int = 0
)
const TypeSleepDiscoveryMode = "E8"

type SleepDiscoveryMode struct {
	*Int
}

func NewSleepDiscoveryMode() *SleepDiscoveryMode {

	char := NewInt(TypeSleepDiscoveryMode)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &SleepDiscoveryMode{char}
}
