package addons

import (
	"fmt"
	addon "gitee.com/liu_guilin/gateway-addon-golang"
)

var PropertyChangedChan chan *addon.Property

type PropertyProxy struct {
	deviceProxy *DeviceProxy
	*addon.Property
}

func NewPropertyProxy(devProxy *DeviceProxy, property *addon.Property) *PropertyProxy {
	properProxy := &PropertyProxy{
		deviceProxy: devProxy,
		Property:    property,
	}
	property.DoPropertyChanged = properProxy.OnPropertyChanged()
	return properProxy
}

func (proxy *PropertyProxy) OnPropertyChanged() func(p *addon.Property) {
	return func(new *addon.Property) {
		select {
		case PropertyChangedChan <- new:
		default:
			if proxy.Value != new.Value {
				_ = proxy.SetCachedValue(new.Value)
			}
			if proxy.Title != new.Title {
				proxy.Title = new.Title
			}
			if proxy.Type != new.Type {
				proxy.Type = new.Type
			}
			if proxy.Unit != new.Unit {
				proxy.Unit = new.Unit
			}
			if proxy.Description != new.Description {
				proxy.Description = new.Description
			}
			if proxy.ReadOnly != new.ReadOnly {
				proxy.ReadOnly = new.ReadOnly
			}
		}
	}
}

func (proxy *PropertyProxy) setValue(property *addon.Property, value interface{}) (interface{}, error) {

	proxy.deviceProxy.adapter.sendMessage(DeviceSetPropertyCommand, struct {
		AdapterId     string `json:"adapterId"`
		DeviceId      string `json:"deviceId"`
		PropertyName  string `json:"propertyName"`
		PropertyValue interface{}
	}{AdapterId: proxy.deviceProxy.adapter.ID, DeviceId: proxy.deviceProxy.ID, PropertyName: property.Name, PropertyValue: value})
	var p *addon.Property
	p = <-PropertyChangedChan
	if p.Name == property.Name {
		return p.Value, nil
	}
	return value, fmt.Errorf("deviceProxy set proptery err")
}
