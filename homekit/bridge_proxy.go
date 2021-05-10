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

var Bridge *BridgeProxy
var stop func()

// HService 不同的service实现此接口
type HService interface {
	GetService() *service.Service
	GetHCharacteristic(name string) HCharacteristic
	GetID() string
}

type BridgeProxy struct {
	*accessory.Bridge
	things   *models.Things
	services map[string]HService
	info     accessory.Info
	config   hc.Config
}

func NewHomeKitBridge(name, manufacturer, model, storagePath, pin string) *BridgeProxy {

	b := &BridgeProxy{}
	b.info = accessory.Info{
		Name:             name,
		Manufacturer:     manufacturer,
		Model:            model,
		FirmwareRevision: FirmwareRevision,
	}
	b.config = hc.Config{
		Port:        "22345",
		StoragePath: storagePath,
		Pin:         pin,
	}
	b.Bridge = accessory.NewBridge(b.info)
	return b
}

// GetThings get things form things models,add service to bridge
func (br *BridgeProxy) GetThings() {
	mapOfThings := br.things.GetThings()
	for _, thing := range mapOfThings {
		s := NewHomekitService(thing)
		if s == nil {
			continue
		}
		br.addService(s)
	}
}

// bridge append service
func (br *BridgeProxy) addService(s HService) {
	br.services[s.GetID()] = s
	br.Bridge.AddService(s.GetService())
}

func (br *BridgeProxy) updateServices() {
	br.Services = make([]*service.Service, 1)
	br.GetThings()
}

func (br *BridgeProxy) onThingAdded(data []byte) {
	deviceId := json.Get(data, "deviceId").ToString()
	log.Info(deviceId)

}

func (br *BridgeProxy) onThingsModified(data []byte) {

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
	c.OnPropertChanged(value)
}

func (br *BridgeProxy) Start() error {

	t, err := hc.NewIPTransport(br.config, br.Accessory)
	if err != nil {
		return fmt.Errorf("create homekit transport err:%s", err.Error())

	}
	stop = func() {
		<-t.Stop()
	}
	log.Info("homekit bridge start")
	go t.Start()
	br.things.Subscribe(util.ThingAdded, Bridge.onThingAdded)
	plugin.Subscribe(util.PropertyChanged, Bridge.onPropertyChanged)
	br.things.Subscribe(util.ThingModified, Bridge.onThingsModified)
	return nil
}

func (br *BridgeProxy) Stop() {
	br.things.Unsubscribe(util.ThingAdded, Bridge.onThingAdded)
	plugin.Unsubscribe(util.PropertyChanged, Bridge.onPropertyChanged)
	br.things.Unsubscribe(util.ThingModified, Bridge.onThingsModified)
	stop()
}
