package properties

import "fmt"

type ColorProperty struct {
	*StringProperty
}

func NewColorProperty(p PropertyDescription) *ColorProperty {
	prop := &ColorProperty{}
	p.Name = "color"
	p.AtType = "ColorProperty"
	prop.StringProperty = NewStringProperty(p)
	return prop
}

func (c *ColorProperty) SetValue(v string) {
	fmt.Print("set value func not implemented")
}
