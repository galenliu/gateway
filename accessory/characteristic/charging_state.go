//THis File is AUTO-GENERATED
package characteristic

const (
	ChargingStateCharging      int = 1
	ChargingStateNotChargeable int = 2
	ChargingStateNotCharging   int = 0
)
const TypeChargingState = "8F"

type ChargingState struct {
	*Int
}

func NewChargingState() *ChargingState {

	char := NewInt(TypeChargingState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &ChargingState{char}
}
