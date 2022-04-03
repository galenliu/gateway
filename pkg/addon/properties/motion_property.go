package properties

type MotionProperty struct {
	*Boolean
}

func NewMotionProperty(value bool, opts ...Option) *MotionProperty {
	p := &MotionProperty{}
	opts = append(opts, WithTitle("Motion"))
	p.Boolean = NewBoolean(BooleanPropertyDescription{
		Name:     "motion",
		AtType:   TypeMotionProperty,
		ReadOnly: true,
		Value:    value,
	}, opts...)
	return p
}
