//THis File is AUTO-GENERATED
package characteristic

const (
	CurrentHeaterCoolerStateCooling  int = 3
	CurrentHeaterCoolerStateHeating  int = 2
	CurrentHeaterCoolerStateIdle     int = 1
	CurrentHeaterCoolerStateInactive int = 0
)
const TypeCurrentHeaterCoolerState = "B1"

type CurrentHeaterCoolerState struct {
	*Int
}

func NewCurrentHeaterCoolerState() *CurrentHeaterCoolerState {

	char := NewInt(TypeCurrentHeaterCoolerState)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &CurrentHeaterCoolerState{char}
}
