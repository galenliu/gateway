package yeelight

import (
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/addon/proxy"
)

type On struct {
	*properties.OnOffProperty
}

func NewOn(p properties.PropertyDescription) proxy.PropertyProxy {
	return proxy.NewOnOffPropertyInstance(&On{
		properties.NewOnOffProperty(p),
	})
}

func (on *On) SetValue(b bool) {

}
