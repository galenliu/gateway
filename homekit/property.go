package homekit

import (
	"github.com/brutella/hc/characteristic"
	"github.com/galenliu/gateway-addon/properties"
	"github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/server/models"
	"github.com/xiam/to"
	"log"
)

type HCharacteristic interface {
	GetName() string
	SetValue(value interface{})
	GetCharacteristic() *characteristic.Characteristic
	SetChangedFunc(func(interface{}))
	OnPropertChanged(value interface{})
}

// CharacteristicProxy 需要实现characteristic的SetValue和OnValueRemoteUpdate二个业务。
type CharacteristicProxy struct {
	*characteristic.Characteristic
	*models.Property
	onChangedFunc func(interface{})
}

func (c *CharacteristicProxy) GetCharacteristic() *characteristic.Characteristic {
	return c.Characteristic
}

func NewCharacteristicProxy(property *models.Property) HCharacteristic {
	hc := &CharacteristicProxy{}
	switch property.AtType {
	case properties.TypeOnOffProperty:
		c := characteristic.NewOn()
		hc.Characteristic = c.Characteristic
		c.OnValueRemoteUpdate(func(b bool) {
			hc.SetValue(b)
		})
		hc.SetChangedFunc(func(i interface{}) {
			b := to.Bool(i)
			c.SetValue(b)
		})
		return hc
	default:
		return nil
	}
}

func (c *CharacteristicProxy) SetChangedFunc(f func(interface{})) {
	c.onChangedFunc = f
}

func (c *CharacteristicProxy) OnPropertChanged(value interface{}) {
	if c.onChangedFunc != nil {
		c.onChangedFunc(value)
	}
}

func (c *CharacteristicProxy) SetValue(value interface{}) {
	_, err := plugin.SetProperty(c.GetThingId(), c.GetName(), value)
	if err != nil {
		log.Println(err.Error())
	}

}
