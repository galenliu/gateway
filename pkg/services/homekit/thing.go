package homekit

import (
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	addon "github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/server/models"
	"github.com/xiam/to"
)

type ThingProxy interface {
	GetMapOfServiceProxy() map[string]ServiceProxy
	GetPropertyProxy(string2 string) PropertyProxy
	GetID() string
}

type _thing struct {
	*models.Thing
	servicesProxy map[string]ServiceProxy
}

func (t *_thing) GetMapOfServiceProxy() map[string]ServiceProxy {
	return t.servicesProxy
}

func NewThingProxy(th *models.Thing) ThingProxy {
	t := &_thing{}
	t.Thing = th
	t.servicesProxy = make(map[string]ServiceProxy)
	switch th.GetSelectedCapability() {
	case addon.Light:
		sev := service.NewLightbulb()

		_ser := &_service{}
		_ser.Service = sev.Service

		t.servicesProxy[sev.Type] = _ser

		for name, prop := range th.Properties {
			switch prop.GetAtType() {
			case addon.OnOffSwitch:
				sev.On = characteristic.NewOn()
				sev.On.OnValueRemoteUpdate(func(b bool) {
					t.propertyChangedHandler(name, b)
				})
				_ser.propertiesProxy[name] = NewPropertyProxy(th.GetSelectedCapability(), prop, func(value interface{}) {
					sev.On.SetValue(to.Bool(value))
				}, nil)
			}
		}
	}

	return t
}

func (t *_thing) GetMapServiceProxy() map[string]ServiceProxy {
	if len(t.servicesProxy) > 0 {
		return t.servicesProxy
	}
	return nil
}

// GetPropertyProxy return propertyProxy flow the pName
func (t *_thing) GetPropertyProxy(pName string) PropertyProxy {

	for _, serviceProxy := range t.GetMapOfServiceProxy() {
		for name, p := range serviceProxy.GetMapOfPropertyProxy() {
			if name == pName {
				return p
			}
		}
	}
	return nil
}

func (t *_thing) propertyChangedHandler(name string, value interface{}) {
	_, err := plugin.SetProperty(t.Thing.ID, name, value)
	if err != nil {
		return
	}
}

//func (t *_thing) addProperty(ser *service.Service) {
//	for name, prop := range t.Properties {
//		_, ok := t.propertiesProxy[name]
//		if ok {
//			continue
//		}
//		p := NewPropertyProxy(t.GetSelectedCapability(), prop, nil, func(value interface{}) {
//			t.propertyChangedHandler(name, value)
//		})
//		t.propertiesProxy[name] = p
//		ser.AddCharacteristic(p.GetCharacteristic())
//	}
//}
