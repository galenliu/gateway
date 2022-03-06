package properties

type LevelProperty struct {
	*NumberProperty
}

func NewLevelProperty(value Number, min, max Number, opts ...Option) *LevelProperty {

	l := &LevelProperty{}
	l.NumberProperty = NewNumberProperty(NumberPropertyDescription{
		Name:     "level",
		AtType:   TypeLevelProperty,
		Minimum:  min,
		Maximum:  max,
		ReadOnly: false,
		Value:    value,
	}, opts...)
	return l
}
