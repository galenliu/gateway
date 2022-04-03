package properties

type VoltageProperty struct {
	*NumberProperty
}

func NewVoltageProperty(value Number, opts ...Option) *InstantaneousPowerProperty {
	b := &InstantaneousPowerProperty{}
	opts = append(opts, WithTitle("Frequency"), WithUnit(UnitHertz))
	b.NumberProperty = NewNumberProperty(NumberPropertyDescription{
		Name:     "frequency",
		ReadOnly: true,
		AtType:   TypeVoltageProperty,
		Minimum:  0,
		Value:    value,
	}, opts...)
	return b
}
