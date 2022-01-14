package yeelight

import "github.com/galenliu/gateway/pkg/addon/properties"

type On struct {
	*properties.OnOffProperty
}

func NewOn(p properties.PropertyDescription) *On {
	return &On{
		properties.NewOnOffProperty(p),
	}
}

func (on *On) SetValue(b bool) {

}
