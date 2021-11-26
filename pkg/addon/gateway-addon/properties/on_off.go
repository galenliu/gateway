package properties

const TypeOnOffProperty = "OnOffProperty"

type OnOffProperty struct {
	*BooleanProperty
}

func NewOnOffProperty() *OnOffProperty {
	p := &OnOffProperty{}
	p.Type = TypeBoolean
	p.Name = On
	//p.SetValue(false)
	return p
}
