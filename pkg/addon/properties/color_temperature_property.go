package properties

type ColorTemperaturePropertyDescriptor struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Value       Integer `json:"value,omitempty"`
}

type ColorTemperatureProperty struct {
	*IntegerProperty
}

func NewColorTemperatureProperty(value Integer, opts ...Option) *ColorTemperatureProperty {
	colorTemperature := &ColorTemperatureProperty{}
	colorTemperature.IntegerProperty = NewIntegerProperty(IntegerPropertyDescription{
		Name:    "ct",
		AtType:  TypeColorTemperatureProperty,
		Unit:    UnitKelvin,
		Minimum: 2000,
		Maximum: 8000,
		Value:   value,
	}, opts...)
	return colorTemperature
}
