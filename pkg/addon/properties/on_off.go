package properties

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/devices"
	"github.com/galenliu/gateway/pkg/addon/proxy"
)

const On = "on"

type BooleanValuer interface {
	devices.PropertyEntity
	SetValue(b bool)
}

type OnOffProperty struct {
	BooleanValuer
}

func NewOnOffProperty(prop PropertyDescription) *OnOffProperty {
	p := &OnOffProperty{}
	p.BooleanValuer = NewBooleanProperty(prop)
	return p
}

func (p *OnOffProperty) SetValue(a any) {
	fmt.Printf("property: %s SetValue func not implemented", p.GetName())
}

func (p *OnOffProperty) SetHandler(handler proxy.DeviceProxy) {

}
