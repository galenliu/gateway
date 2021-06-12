package homekit

import (
	"fmt"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/plugin"
	"github.com/galenliu/gateway/server/models"
	json "github.com/json-iterator/go"
)

var _bridge *BridgeProxy
var stop func()

// HService 不同的service实现此接口
type HService interface {

}

type BridgeProxy struct {
	*accessory.Bridge
	things   *models.Things
	services map[string]HService
	info     accessory.Info
	config   hc.Config
}

func NewHomeKitBridge(name, manufacturer, model, storagePath, pin string) *BridgeProxy {

	_bridge = &BridgeProxy{}
	_bridge.things = models.NewThings()
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
	_bridge.services = make(map[string]HService)

	light := service.NewLightbulb()
	light.On.OnValueRemoteUpdate(func(b bool) {
		log.Info("home kit light: %b", b)
	})
	_bridge.AddService(light.Service)

	return _bridge
}

// GetThings get things form things models,add service to bridge
func (br *BridgeProxy) GetThings() {
	mapOfThings := br.things.GetThings()
	for _, thing := range mapOfThings {
		s := NewThingProxy(thing)
		if s == nil {
			continue
		}
		br.addService(s)
	}
}

// bridge append service
func (br *BridgeProxy) addService(s ThingProxy) {
	br.services[s.GetID()] = s
	br.AddService(s.GetService())
}

func (br *BridgeProxy) updateServices() {
	br.Services = make([]*service.Service, 1)
	br.GetThings()
}

func (br *BridgeProxy) onThingAdded(thing *models.Thing) {
	br.GetThings()
	log.Info("added thing %v", thing)
}

func (br *BridgeProxy) onThingsModified(data []byte) {
	log.Error(string(data))
}

func (br *BridgeProxy) onPropertyChanged(data []byte) {
	deviceId := json.Get(data, "deviceId").ToString()
	propName := json.Get(data, "name").ToString()
	value := json.Get(data, "value").GetInterface()
	if deviceId == "" || propName == "" || value == nil {
		return
	}
	ser := br.services[deviceId]
	if ser == nil {
		return
	}
	c := ser.GetHCharacteristic(propName)
	if c == nil {
		return
	}
	c.OnPropertyChanged(value)
}

func (br *BridgeProxy) Start() error {

	t, err := hc.NewIPTransport(br.config, br.Bridge.Accessory)
	if err != nil {
		return fmt.Errorf("create homekit transport err:%s", err.Error())

	}
	stop = func() {
		<-t.Stop()
	}
	log.Info("homekit bridge start")
	go t.Start()
	br.GetThings()
	br.things.Subscribe(util.ThingAdded, _bridge.onThingAdded)
	plugin.Subscribe(util.PropertyChanged, _bridge.onPropertyChanged)
	br.things.Subscribe(util.ThingModified, _bridge.onThingsModified)
	return nil
}

func (br *BridgeProxy) Stop() {
	br.things.Unsubscribe(util.ThingAdded, _bridge.onThingAdded)
	plugin.Unsubscribe(util.PropertyChanged, _bridge.onPropertyChanged)
	br.things.Unsubscribe(util.ThingModified, _bridge.onThingsModified)
	stop()
}
