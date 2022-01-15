package properties

import (
	"fmt"
)

type OnOffProperty struct {
	*BooleanProperty
}

func NewOnOffProperty(prop PropertyDescription) *OnOffProperty {
	p := &OnOffProperty{}
	var name = "on"
	prop.Name = &name
	atType := "OnOffProperty"
	prop.AtType = &atType
	p.BooleanProperty = NewBooleanProperty(prop)
	return p
}

func (p *OnOffProperty) SetValue(a any) {
	fmt.Printf("property: %s SetValue func not implemented", p.GetName())
}

func (p *OnOffProperty) TurnOn() {
	fmt.Printf("property: %s TurnOn func not implemented", p.GetName())
}

func (p *OnOffProperty) TurnOff() {
	fmt.Printf("property: %s TurnOff func not implemented", p.GetName())
}
