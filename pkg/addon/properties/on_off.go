package properties

type OnOff interface {
	BooleanEntity
}

type OnOffProperty struct {
	*BooleanProperty
}

func NewOnOffProperty(prop PropertyDescription) *OnOffProperty {
	p := &OnOffProperty{}
	prop.Name = "on"
	prop.AtType = "OnOffProperty"
	p.BooleanProperty = NewBooleanProperty(prop)
	return p
}
