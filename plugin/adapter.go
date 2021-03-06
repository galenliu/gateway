package plugin

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway/plugin/internal"
	"sync"
)

type pairingFunc func(ctx context.Context, cancelFunc func())

type managerProxy interface {
	handleDeviceAdded(device *addon.Device)
}

type Adapter struct {
	id             string
	name           string
	pluginId       string
	plugin         *Plugin
	looker         *sync.Mutex
	isPairing      bool
	onPairingFunc  pairingFunc
	devices        map[string]addon.IDevice
	pairingContext context.Context
	manifest       interface{}
	packageName    string
	manager        managerProxy
}

func NewAdapter(manager managerProxy, name, adapterId, pluginId, packageName string) *Adapter {
	proxy := &Adapter{}
	proxy.id = adapterId
	proxy.name = name
	proxy.packageName = packageName
	proxy.pluginId = pluginId
	proxy.devices = make(map[string]addon.IDevice)
	proxy.looker = new(sync.Mutex)
	proxy.manager = manager
	return proxy
}

func (adapter *Adapter) pairing(timeout float64) {
	logging.Info(fmt.Sprintf("adapter: %s start pairing", adapter.id))
	data := make(map[string]interface{})
	data["timeout"] = timeout
	adapter.Send(internal.AdapterStartPairingCommand, data)
}

func (adapter *Adapter) cancelPairing() {
	logging.Info(fmt.Sprintf("adapter: %s execute pairing", adapter.id))
	data := make(map[string]interface{})
	adapter.Send(internal.AdapterCancelPairingCommand, data)
}

func (adapter *Adapter) removeThing(device addon.IDevice) {
	logging.Info(fmt.Sprintf("adapter delete thing Id: %v", device.GetID()))
	data := make(map[string]interface{})
	data["deviceId"] = device.GetID()
	adapter.Send(internal.AdapterRemoveDeviceRequest, data)

}

func (adapter *Adapter) cancelRemoveThing(deviceId string) {
	logging.Info(fmt.Sprintf("adapter: %s execute pairing", adapter.id))
	data := make(map[string]interface{})
	data["deviceId"] = deviceId
	adapter.Send(internal.AdapterCancelRemoveDeviceCommand, data)
}

func (adapter *Adapter) getManager() *manager {
	return adapter.plugin.pluginServer.manager
}

func (adapter *Adapter) Send(messageType int, data map[string]interface{}) {
	data["adapterId"] = adapter.id
	adapter.plugin.send(messageType, data)
}

func (adapter *Adapter) getDevice(deviceId string) addon.IDevice {
	device, ok := adapter.devices[deviceId]
	if ok {
		return device
	}
	return nil
}

func (adapter *Adapter) handleDeviceRemoved(device addon.IDevice) {
	delete(adapter.devices, device.GetID())

}

func (adapter *Adapter) handleDeviceAdded(device *addon.Device) {
	adapter.devices[device.GetID()] = device
	adapter.manager.handleDeviceAdded(device)
}
