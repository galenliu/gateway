//THis File is AUTO-GENERATED
package characteristic

const (
	ProgrammableSwitchEventDoublePress int = 1
	ProgrammableSwitchEventLongPress   int = 2
	ProgrammableSwitchEventSinglePress int = 0
)
const TypeProgrammableSwitchEvent = "73"

type ProgrammableSwitchEvent struct {
	*Int
}

func NewProgrammableSwitchEvent() *ProgrammableSwitchEvent {

	char := NewInt(TypeProgrammableSwitchEvent)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &ProgrammableSwitchEvent{char}
}
