package properties

//type OnOffPropertyDescription struct {
//	Description string `json:"description"`
//	Value       bool   `json:"value"`
//}

type OnOff interface {
	BooleanEntity
}

type OnOffProperty struct {
	*BooleanProperty
}

func NewOnOffProperty(value bool, args ...string) *OnOffProperty {

	p := &OnOffProperty{}
	p.BooleanProperty = NewBooleanProperty(BooleanPropertyDescription{
		Name:   "on",
		AtType: TypeOnOffProperty,
		Value:  value,
	})
	return p
}
