//THis File is AUTO-GENERATED
package characteristic

const (
	PictureModeCalibrated     int = 2
	PictureModeCalibratedDark int = 3
	PictureModeComputer       int = 6
	PictureModeCustom         int = 7
	PictureModeGame           int = 5
	PictureModeOther          int = 0
	PictureModeStandard       int = 1
	PictureModeVivid          int = 4
)
const TypePictureMode = "E2"

type PictureMode struct {
	*Int
}

func NewPictureMode() *PictureMode {

	char := NewInt(TypePictureMode)
	char.Format = FormatUint8
	char.Perms = []string{PermRead, PermWrite, PermEvents}
	char.SetMinValue(0)
	char.SetMaxValue(13)

	char.SetValue(0)
	return &PictureMode{char}
}
