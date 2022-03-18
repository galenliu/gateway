package properties

type ColorModePropertyEnum = string

const ColorModePropertyEnumColor = "color"
const ColorModePropertyEnumTemperature = "temperature"

type ColorModeProperty struct {
	*StringProperty
}

func NewColorModeProperty(value ColorModePropertyEnum, opts ...Option) *ColorModeProperty {
	p := &ColorModeProperty{}
	p.StringProperty = NewStringProperty(StringPropertyDescription{
		Name:     "color_mode",
		AtType:   TypeColorModeProperty,
		Title:    "color mode",
		Type:     TypeString,
		Enum:     []string{ColorModePropertyEnumColor, ColorModePropertyEnumTemperature},
		ReadOnly: false,
		Value:    value,
	}, opts...)
	return p
}
