package properties

type TargetTemperatureProperty struct {
	*NumberProperty
}

func NewTargetTemperatureProperty(value Number, opts ...Option) *TargetTemperatureProperty {

	b := &TargetTemperatureProperty{}
	opts = append(opts, WithUnit(UnitArcDegrees))
	b.NumberProperty = NewNumberProperty(NumberPropertyDescription{
		Name:       "targetTemperature",
		AtType:     TypeTargetTemperatureProperty,
		Minimum:    10,
		Maximum:    38,
		MultipleOf: 0.1,
		Value:      value,
	}, opts...)
	return b
}
