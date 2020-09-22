//THis File is AUTO-GENERATED
package characteristic

const (
	RemoteKeyArrowDown     int = 5
	RemoteKeyArrowLeft     int = 6
	RemoteKeyArrowRight    int = 7
	RemoteKeyArrowUp       int = 4
	RemoteKeyBack          int = 9
	RemoteKeyExit          int = 10
	RemoteKeyFastForward   int = 1
	RemoteKeyInformation   int = 15
	RemoteKeyNextTrack     int = 2
	RemoteKeyPlayPause     int = 11
	RemoteKeyPreviousTrack int = 3
	RemoteKeyRewind        int = 0
	RemoteKeySelect        int = 8
)
const TypeRemoteKey = "E1"

type RemoteKey struct {
	*Int
}

func NewRemoteKey() *RemoteKey {

	char := NewInt(TypeRemoteKey)
	char.Format = FormatUint8
	char.Perms = []string{PermWrite}
	char.SetMinValue(0)
	char.SetMaxValue(16)

	char.SetValue(0)
	return &RemoteKey{char}
}
