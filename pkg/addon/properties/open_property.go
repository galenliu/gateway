package properties

type OpenProperty struct {
	*Boolean
}

func NewOpenProperty(value bool, opts ...Option) *MotionProperty {
	p := &MotionProperty{}
	p.Boolean = NewBoolean(BooleanPropertyDescription{
		Name:   "open",
		AtType: TypeOpenProperty,
		Title:  "Open",
		Value:  value,
	}, opts...)
	return p
}
