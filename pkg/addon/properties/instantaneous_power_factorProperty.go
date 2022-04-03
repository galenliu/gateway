package properties

type InstantaneousPowerFactorProperty struct {
	*NumberProperty
}

func NewInstantaneousPowerFactorProperty(value Number, opts ...Option) *InstantaneousPowerProperty {
	b := &InstantaneousPowerProperty{}
	opts = append(opts, WithTitle("Power"), WithUnit(UnitWatt))
	b.NumberProperty = NewNumberProperty(NumberPropertyDescription{
		Name:    "power",
		AtType:  TypeInstantaneousPowerFactorProperty,
		Minimum: -1,
		Maximum: 1,
		Value:   value,
	}, opts...)
	return b
}
