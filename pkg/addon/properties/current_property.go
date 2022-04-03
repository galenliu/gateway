package properties

type CurrentProperty struct {
	*NumberProperty
}

func NewCurrentProperty(value Number, opts ...Option) *InstantaneousPowerProperty {
	b := &InstantaneousPowerProperty{}
	opts = append(opts, WithTitle("Current"), WithUnit(UnitAmpere))
	b.NumberProperty = NewNumberProperty(NumberPropertyDescription{
		Name:   "current",
		AtType: TypeCurrentProperty,
		Value:  value,
	}, opts...)
	return b
}
