package homekit

import (
	"fmt"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/server/models"
	json "github.com/json-iterator/go"
)

var _bridge *bridge
var stop func()

// HService 不同的service实现此接口
type HService interface {
}

type BridgeProxy interface {
	Start() error
	Stop()
}

type bridge struct {
	*accessory.Bridge
	thingsProxy map[string]ThingProxy
	info        accessory.Info
	config      hc.Config
	stopped     <-chan struct{}
}

func NewHomeKitBridge(name, manufacturer, model, storagePath, pin string) BridgeProxy {

	_bridge = &bridge{}
	_bridge.info = accessory.Info{
		Name:             name,
		Manufacturer:     manufacturer,
		Model:            model,
		FirmwareRevision: FirmwareRevision,
	}
	_bridge.config = hc.Config{
		Port:        "22345",
		StoragePath: storagePath,
		Pin:         pin,
	}
	_bridge.Bridge = accessory.NewBridge(_bridge.info)
	_bridge.stopped = make(chan struct{})
	_bridge.updateThingsProxy()
	return _bridge

}

// GetThings get things form things models,add service to bridge
//func (br *bridge) GetThings() {
//	mapOfThings := br.things.GetThings()
//	for _, thing := range mapOfThings {
//		s := NewThingProxy(thing)
//		if s == nil {
//			continue
//		}
//		br.addService(s)
//	}
//}

//func (br *bridge) updateServices() {
//	br.Services = make([]*service.Service, 1)
//	br.GetThings()
//}

func (br *bridge) Start() error {
	t, err := hc.NewIPTransport(br.config, br.Bridge.Accessory)
	if err != nil {
		return fmt.Errorf("create homekit transport err:%s", err.Error())
	}
	br.stopped = t.Stop()
	logging.Info("homekit bridge start")
	models.NewThingsOnce().Subscribe(util.ThingAdded, _bridge.onThingAdded)
	plugin.Subscribe(util.PropertyChanged, _bridge.onPropertyChanged)
	models.NewThingsOnce().Subscribe(util.ThingModified, _bridge.onThingsModified)
	go t.Start()
	return nil
}

func (br *bridge) Stop() {
	select {
	case <-br.stopped:
		logging.Info("HomeKit Bridge Stop")
	}
	models.NewThingsOnce().Unsubscribe(util.ThingAdded, _bridge.onThingAdded)
	plugin.Unsubscribe(util.PropertyChanged, _bridge.onPropertyChanged)
	models.NewThingsOnce().Unsubscribe(util.ThingModified, _bridge.onThingsModified)
}

func (br *bridge) onThingAdded(thing *models.Thing) {
	br.restart()
}

func (br *bridge) onThingsModified(data []byte) {
	br.restart()
}

func (br *bridge) onPropertyChanged(data []byte) {
	deviceId := json.Get(data, "deviceId").ToString()
	propName := json.Get(data, "name").ToString()
	value := json.Get(data, "value").GetInterface()
	if deviceId == "" || propName == "" || value == nil {
		return
	}
	_thingProxy, ok := br.thingsProxy[deviceId]
	if !ok {
		logging.Error("device id:&s err", deviceId)
		return
	}
	_propProxy := _thingProxy.GetPropertyProxy(propName)
	_propProxy.SetValue(value)
}

func (br *bridge) restart() {
	br.Stop()
	br.updateThingsProxy()
	err := br.Start()
	if err != nil {
		logging.Error("HomeKit Bridge err: %s", err.Error())
	}
}

func (br *bridge) updateThingsProxy() {
	br.thingsProxy = make(map[string]ThingProxy)
	for name, thing := range models.NewThingsOnce().GetMapOfThings() {
		_proxy := NewThingProxy(thing)
		if _proxy != nil {
			br.thingsProxy[name] = _proxy
		}
	}
	for _, propertyProxy := range br.thingsProxy {
		propertyProxy.GetMapOfServiceProxy()
	}
	br.Bridge.AddService(nil)
}
