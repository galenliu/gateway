package properties

type TemperatureProperty struct {
	*NumberProperty
}

func NewTemperatureProperty(value Number, opts ...Option) *TemperatureProperty {
	b := &TemperatureProperty{}
	opts = append(opts, WithTitle("Temperature"), WithUnit(UnitArcDegrees))
	b.NumberProperty = NewNumberProperty(NumberPropertyDescription{
		Name:     "temperature",
		AtType:   TypeTemperatureProperty,
		ReadOnly: true,
		Minimum:  -100,
		Maximum:  100,
		Value:    value,
	}, opts...)
	return b
}
