package properties

type HeatingCoolingProperty struct {
	*StringProperty
}

func NewHeatingCoolingProperty(value ThermostatModeEnum, opts ...Option) *HeatingCoolingProperty {
	b := &HeatingCoolingProperty{}
	b.StringProperty = NewStringProperty(StringPropertyDescription{
		Name:   "heatingCooling",
		Enum:   []string{"off", "heat", "cool"},
		Title:  "HeatingCooling",
		Unit:   UnitDegreeCelsius,
		AtType: TypeHeatingCoolingProperty,
		Value:  string(value),
	}, opts...)
	return b
}
