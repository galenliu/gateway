package properties

type FrequencyProperty struct {
	*NumberProperty
}

func NewFrequencyProperty(value Number, opts ...Option) *InstantaneousPowerProperty {
	b := &InstantaneousPowerProperty{}
	opts = append(opts, WithTitle("Frequency"), WithUnit(UnitVolt))
	b.NumberProperty = NewNumberProperty(NumberPropertyDescription{
		Name:     "frequency",
		ReadOnly: true,
		AtType:   TypeFrequencyProperty,
		Value:    value,
	}, opts...)
	return b
}
