package properties

type LeakProperty struct {
	*Boolean
}

func NewLeakProperty(value bool, opts ...Option) *LeakProperty {
	p := &LeakProperty{}
	p.Boolean = NewBoolean(BooleanPropertyDescription{
		Name:   "leak",
		AtType: TypeLeakProperty,
		Title:  "Leak",
		Value:  value,
	}, opts...)
	return p
}
