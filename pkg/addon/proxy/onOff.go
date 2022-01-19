package proxy

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/properties"
)

type OnOffPropertyInstance interface {
	properties.Entity
	TurnOn() error
	TurnOff() error
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
			err := on.OnOffPropertyInstance.TurnOn()
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
		} else {
			err := on.OnOffPropertyInstance.TurnOff()
			if err != nil {
				fmt.Printf(err.Error())
				return
			}
		}
		on.SetCachedValueAndNotify(b)
	}
}
