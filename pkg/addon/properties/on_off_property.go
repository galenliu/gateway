package properties

//type OnOffPropertyDescription struct {
//	Description string `json:"description"`
//	Value       bool   `json:"value"`
//}

type OnOff interface {
	BooleanEntity
}

type OnOffProperty struct {
	*Boolean
}

func NewOnOffProperty(value bool, args ...string) *OnOffProperty {

	p := &OnOffProperty{}
	p.Boolean = NewBoolean(BooleanPropertyDescription{
		Name:   "on",
		AtType: TypeOnOffProperty,
		Value:  value,
	})
	return p
}
