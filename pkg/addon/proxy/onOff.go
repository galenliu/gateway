package proxy

import "github.com/galenliu/gateway/pkg/addon/devices"

type OnOffPropertyInstance interface {
	devices.PropertyEntity
	TurnOn()
	TurnOff()
}

type OnOfProperty struct {
	OnOffPropertyInstance
}

func NewOnOff(onOff OnOffPropertyInstance) *OnOfProperty {
	return &OnOfProperty{onOff}
}

func (on *OnOfProperty) SetValue(a any) {
	b, ok := a.(bool)
	if ok {
		if b {
			on.OnOffPropertyInstance.TurnOn()
		} else {
			on.OnOffPropertyInstance.TurnOff()
		}
	}
}
