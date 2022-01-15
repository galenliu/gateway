package properties

import "fmt"

type ColorProperty struct {
	*StringProperty
}

func NewColorProperty(p PropertyDescription) *ColorProperty {
	prop := &ColorProperty{}
	var name = "color"
	p.Name = &name
	var atType = "ColorProperty"
	p.AtType = &atType
	prop.StringProperty = NewStringProperty(p)
	return prop
}

func (c *ColorProperty) SetValue(v string) {
	fmt.Print("set value func not implemented")
}
