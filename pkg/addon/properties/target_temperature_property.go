package properties

type TargetTemperatureProperty struct {
	*NumberProperty
}

func NewTargetTemperatureProperty(value Number, opts ...Option) *TargetTemperatureProperty {

	b := &TargetTemperatureProperty{}
	b.NumberProperty = NewNumberProperty(NumberPropertyDescription{
		Name:       "targetTemperature",
		Unit:       UnitDegreeCelsius,
		Title:      "TargetTemperature",
		AtType:     TypeTargetTemperatureProperty,
		Minimum:    10,
		Maximum:    38,
		MultipleOf: 0.1,
		Value:      value,
	}, opts...)
	return b
}
