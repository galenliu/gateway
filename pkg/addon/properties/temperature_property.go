package properties

type TemperatureProperty struct {
	*NumberProperty
}

func NewTemperatureProperty(value Number, opts ...Option) *TemperatureProperty {
	b := &TemperatureProperty{}
	opts = append(opts, WithTitle("Temperature"))
	b.NumberProperty = NewNumberProperty(NumberPropertyDescription{
		Name:       "temperature",
		Unit:       UnitDegreeCelsius,
		AtType:     TypeTemperatureProperty,
		ReadOnly:   true,
		Title:      "Temperature",
		Minimum:    10,
		Maximum:    38,
		MultipleOf: 0.1,
		Value:      value,
	}, opts...)
	return b
}
