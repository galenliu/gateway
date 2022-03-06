package properties

type ColorTemperaturePropertyDescriptor struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Value       Integer `json:"value,omitempty"`
}

type ColorTemperatureProperty struct {
	*IntegerProperty
}

func NewColorTemperatureProperty(desc ColorTemperaturePropertyDescriptor, opts ...Option) *ColorTemperatureProperty {
	colorTemperature := &ColorTemperatureProperty{}
	colorTemperature.IntegerProperty = NewIntegerProperty(IntegerPropertyDescription{
		Name:    "ct",
		AtType:  TypeInteger,
		Unit:    UnitKelvin,
		Minimum: 2000,
		Maximum: 8000,
		Value:   desc.Value,
	}, opts...)
	return colorTemperature
}