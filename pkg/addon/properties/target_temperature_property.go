package properties

type TargetTemperatureProperty struct {
	*NumberProperty
}

func NewTargetTemperatureProperty(value Number, opts ...Option) *TargetTemperatureProperty {

	b := &TargetTemperatureProperty{}
	opts = append(opts, WithTitle("TargetTemperature"), WithUnit(UnitArcDegrees))
	b.NumberProperty = NewNumberProperty(NumberPropertyDescription{
		Name:    "targetTemperature",
		AtType:  TypeTargetTemperatureProperty,
		Minimum: -100,
		Maximum: 100,
		Value:   value,
	}, opts...)
	return b
}
