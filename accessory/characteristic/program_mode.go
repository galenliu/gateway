//THis File is AUTO-GENERATED
package characteristic

const (
	ProgramModeNoprogramscheduled          int = 0
	ProgramModeProgramscheduled            int = 1
	ProgramModeProgramscheduled_ManualMode int = 2
)
const TypeProgramMode = "D1"

type ProgramMode struct {
	*Int
}

func NewProgramMode() *ProgramMode {

	char := NewInt(TypeProgramMode)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermEvents}

	char.SetValue(0)
	return &ProgramMode{char}
}
