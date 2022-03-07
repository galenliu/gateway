package properties

type BooleanProperty struct {
	*Boolean
}

func NewBooleanProperty(value bool, args ...string) *BooleanProperty {
	p := &BooleanProperty{}
	p.Boolean = NewBoolean(BooleanPropertyDescription{
		Name:   "on",
		AtType: TypeBooleanProperty,
		Value:  value,
	})
	return p
}
