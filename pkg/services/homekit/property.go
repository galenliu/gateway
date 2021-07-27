package homekit

import (
	"github.com/brutella/hc/characteristic"
	addon "github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway-addon/properties"
	"github.com/galenliu/gateway/server/models/model"
	"net"
)

type PropertyProxy interface {
	GetName() string
	SetValue(value interface{})
	GetCharacteristic() *characteristic.Characteristic
}

type CharacteristicProxy interface {
	onChanged(value interface{})
	setValue(value interface{})
}

// property 需要实现characteristic的SetValue和OnValueRemoteUpdate二个业务。
type property struct {
	*characteristic.Characteristic
	*model.Property
	onChangedFunc func(interface{})
	setValue      func(value interface{})
}



func NewPropertyProxy(typ string, p *model.Property,onUpdateValue func(interface{}), onChanged func(interface{})) PropertyProxy {
	proxy := &property{}
	proxy.Property = p

	switch p.AtType {
	case addon.OnOffSwitch:
		switch typ {
		case addon.Light, addon.OnOffSwitch:
			c := characteristic.NewOn()
			proxy.Characteristic = c.Characteristic
		}
	case properties.TypeColorProperty:
		c := characteristic.NewHue()
		proxy.Characteristic = c.Characteristic

	case properties.TypeBrightnessProperty:
		c := characteristic.NewBrightness()
		proxy.Characteristic = c.Characteristic

	default:
		return nil
	}

	if proxy.Characteristic == nil {
		return nil
	}
	if onChanged != nil && proxy.GetCharacteristic() != nil{
		proxy.GetCharacteristic().OnValueUpdateFromConn(func(conn net.Conn, c *characteristic.Characteristic, newValue, oldValue interface{}) {
			onChanged(newValue)
		})
	}

	if onUpdateValue != nil{
		proxy.setValue = onUpdateValue
	}
	return proxy
}

func (p *property) GetCharacteristic() *characteristic.Characteristic {
	return p.Characteristic
}

func (p *property) SetValue(value interface{}) {
	if p.setValue != nil {
		p.setValue(value)
	}
}
