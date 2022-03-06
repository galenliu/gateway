package properties

type ColorModePropertyEnum = string

const ColorModePropertyEnumColor = "color"
const ColorModePropertyEnumTemperature = "temperature"

type ColorModePropertyDescription struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Value       ColorModePropertyEnum
}

type ColorModeProperty struct {
	*StringProperty
}

func NewColorModeProperty(desc ColorModePropertyDescription) *ColorModeProperty {
	p := &ColorModeProperty{}
	p.StringProperty = NewStringProperty(StringPropertyDescription{
		Name:     "color_mode",
		AtType:   TypeColorModeProperty,
		Type:     TypeString,
		Enum:     []string{ColorModePropertyEnumColor, ColorModePropertyEnumTemperature},
		ReadOnly: false,
		Value:    desc.Value,
	})
	return p
}
