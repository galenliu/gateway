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
	desc := ""
	title := "on/off"
	if len(args) > 0 {
		desc = args[0]
	}
	if len(args) > 1 {
		title = args[1]
	}
	p := &OnOffProperty{}
	p.BooleanProperty = NewBooleanProperty(BooleanPropertyDescription{
		Name:        "on",
		AtType:      TypeOnOffProperty,
		Title:       title,
		Description: desc,
		Value:       value,
	})
	return p
}
