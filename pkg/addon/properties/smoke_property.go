package properties

type SmokeProperty struct {
	*Boolean
}

func NewSmokeProperty(value bool, opts ...Option) *SmokeProperty {
	p := &SmokeProperty{}
	p.Boolean = NewBoolean(BooleanPropertyDescription{
		Name:   "open",
		AtType: TypeSmokeProperty,
		Title:  "Smoke",
		Value:  value,
	}, opts...)
	return p
}
