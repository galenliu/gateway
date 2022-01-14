package properties

import (
	"github.com/galenliu/gateway/pkg/addon/proxy"
	"github.com/galenliu/gateway/pkg/addon/schemas"
)

const On = "on"

type OnOffPropertyEntity interface {
	proxy.PropertyProxy
	Turn(b bool)
}

type OnOffProperty struct {
	*BooleanProperty
}

func NewOnOffProperty(prop PropertyDescription) *OnOffProperty {
	p := &OnOffProperty{}
	p.BooleanProperty = NewBooleanProperty(prop)
	p.Type = TypeBoolean
	p.AtType = schemas.OnOffProperty
	p.Name = On
	return p
}

func (p *OnOffProperty) Turn(b bool) {

}

func (p *OnOffProperty) SetValue(a any) {

}
