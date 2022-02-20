package properties

const TypeColorModeProperty = "ColorModeProperty"

type ColorModeProperty struct {
	*StringProperty
}

func NewColorModeProperty() *ColorModeProperty {
	p := &ColorModeProperty{}
	p.Type = TypeString
	p.AtType = TypeColorModeProperty
	p.Name = ColorModel
	p.Enum = []any{"color", "temperature"}
	return p
}
