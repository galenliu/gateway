package properties

type InstantaneousPowerProperty struct {
	*NumberProperty
}

func NewInstantaneousPowerProperty(value Number, opts ...Option) *InstantaneousPowerProperty {

	b := &InstantaneousPowerProperty{}
	opts = append(opts, WithTitle("Power"), WithUnit(UnitWatt))
	b.NumberProperty = NewNumberProperty(NumberPropertyDescription{
		Name:    "power",
		AtType:  TypeInstantaneousPowerProperty,
		Minimum: 0,
		Value:   value,
	}, opts...)
	return b
}
