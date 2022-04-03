package properties

type OnOff interface {
	BooleanEntity
}

type OnOffProperty struct {
	*Boolean
}

func NewOnOffProperty(value bool, opts ...Option) *OnOffProperty {
	p := &OnOffProperty{}
	p.Boolean = NewBoolean(BooleanPropertyDescription{
		Name:   "on",
		AtType: TypeOnOffProperty,
		Value:  value,
	}, opts...)
	return p
}
