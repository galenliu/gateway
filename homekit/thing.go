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
	GetService() ServiceProxy
	GetID() string
}

type ServiceProxy interface {
}

type thing struct {
	*models.Thing
	ServiceProxy
	propertiesProxy map[string]PropertyProxy
}

func NewThingProxy(th *models.Thing) ThingProxy {
	t := &thing{}
	t.Thing = th
	t.propertiesProxy = make(map[string]PropertyProxy)
	switch th.GetSelectedCapability() {
	case addon.Light:
		sev := service.NewLightbulb()
		t.ServiceProxy = sev
		for name, prop := range th.Properties {
			switch prop.GetAtType() {
			case addon.OnOffSwitch:
				sev.On = characteristic.NewOn()
				sev.On.OnValueRemoteUpdate(func(b bool) {
					t.propertyChangedHandler(name, b)
				})
				t.propertiesProxy[name] = NewPropertyProxy(th.GetSelectedCapability(), prop, func(value interface{}) {
					sev.On.SetValue(to.Bool(value))
				}, nil)
			}
		}
		t.addProperty(sev.Service)
	}

	return t
}

func (t *thing) GetService() ServiceProxy {
	return t.ServiceProxy
}

func (t *thing) propertyChangedHandler(name string, value interface{}) {
	_, err := plugin.SetProperty(t.Thing.ID, name, value)
	if err != nil {
		return
	}
}

func (t *thing) addProperty(ser *service.Service) {
	for name, prop := range t.Properties {
		_, ok := t.propertiesProxy[name]
		if ok {
			continue
		}
		p := NewPropertyProxy(t.GetSelectedCapability(), prop,nil, func(value interface{}) {
			t.propertyChangedHandler(name, value)
		}, )
		t.propertiesProxy[name] = p
		ser.AddCharacteristic(p.GetCharacteristic())
	}
}
