package addon

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/gateway-addon/devices"
	"log"
	"sync"
)

type onPairingFunc func(timeout float64)
type OnCancelPairingFunc func()
type OnDeviceSavedFunc func(deivceId string, device *devices.Device)
type OnSetCredentialsFunc func(deivceId, username, password string)

//type OnSetPinFunc func(deivceId string, devices.P) error

type Adapter struct {
	Devices     map[string]IDevice
	manager     *AddonManager
	locker      *sync.Mutex
	cancelChan  chan struct{}
	Id          string
	name        string
	packageName string
	IsPairing   bool
	verbose     bool
}

func NewAdapter(adapterId, adapterName string) *Adapter {

	adapter := &Adapter{}
	adapter.Id = adapterId
	adapter.name = adapterName
	adapter.locker = new(sync.Mutex)
	adapter.Devices = make(map[string]IDevice)
	adapter.cancelChan = make(chan struct{})
	adapter.verbose = true
	return adapter
}

func (a *Adapter) HandleDeviceAdded(device IDevice) {
	a.Devices[device.GetID()] = device
	a.manager.handleDeviceAdded(device)
}

func (a *Adapter) SendError(message string) {
	data := make(map[string]interface{})
	data[Aid] = a.GetID()
	data["message"] = message
	a.manager.send(PluginErrorNotification, data)
}

//func (a *Adapter) ConnectedNotify(device *Device, connected bool) {
//	a.manager.sendConnectedStateNotification(device, connected)
//}

//向前端UI发送提示
func (a *Adapter) SendPairingPrompt(prompt, url string, did string) {
	data := make(map[string]interface{})
	data[Aid] = a.GetID()
	data["prompt"] = prompt
	data[Did] = did
	if url != "" {
		data["url"] = url
	}
	a.manager.send(AdapterPairingPromptNotification, data)
}

func (a *Adapter) SendUnpairingPrompt(prompt, url string, did string) {
	data := make(map[string]interface{})
	data[Aid] = a.GetID()
	data["prompt"] = prompt
	if did != "" {
		data[Did] = did
	}
	if url != "" {
		data["url"] = url
	}
	a.manager.send(AdapterUnpairingPromptNotification, data)
}

func (a *Adapter) Send(mt int, data map[string]interface{}) {
	a.manager.send(mt, data)
}

func (a *Adapter) StartPairing(timeout float64) {
	if a.verbose {
		log.Printf("adapter:(%s)- StartPairing() not implemented", a.GetID())
	}
}

func (a *Adapter) CancelPairing() {
	if a.verbose {
		log.Printf("adapter:(%s)- CancelPairing() not implemented", a.GetID())
	}
}

func (a *Adapter) GetID() string {
	return a.Id
}

func (a *Adapter) GetName() string {
	if a.name == "" {
		return a.Id
	}
	return a.name
}

func (a *Adapter) Unload() {
	if a.verbose {
		log.Printf("adapter:(%s)- unloaded ", a.GetID())
	}
}

func (a *Adapter) HandleDeviceSaved(device IDevice) {
	if a.verbose {
		log.Printf("adapter:(%s)- HandleDeviceSaved() not implemented", a.GetID())
	}
}

func (a *Adapter) HandleDeviceRemoved(device IDevice) {
	delete(a.Devices, device.GetID())
}

func (a *Adapter) getDevice(id string) IDevice {
	return a.Devices[id]
}

func (a *Adapter) close() {
	fmt.Print("do some thing while a close")
	a.manager.close()
}

func (a *Adapter) Running() bool {
	return a.manager.running
}

func (a *Adapter) setManager(m *AddonManager) {
	a.manager = m
}

func (a *Adapter) SetPin(deviceId string, pin interface{}) {

}

func (a *Adapter) SetCredentials(deviceId, username, password string) {

}
