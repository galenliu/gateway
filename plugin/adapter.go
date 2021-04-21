package plugin

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway/pkg/log"
	"sync"
)

type pairingFunc func(ctx context.Context, cancelFunc func())

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
	manager        *AddonManager
}

func NewAdapter(manager *AddonManager, name, adapterId, pluginId, packageName string) *Adapter {
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
	log.Info(fmt.Sprintf("adapter: %s start pairing", adapter.id))
	data := make(map[string]interface{})
	data["timeout"] = timeout
	adapter.send(AdapterStartPairingCommand, data)
}

func (adapter *Adapter) cancelPairing() {
	log.Info(fmt.Sprintf("adapter: %s execute pairing", adapter.id))
	data := make(map[string]interface{})
	adapter.send(AdapterCancelPairingCommand, data)
}

func (adapter *Adapter) removeThing(device addon.IDevice) {
	log.Info(fmt.Sprintf("adapter delete thing Id: %v", device.GetID()))
	data := make(map[string]interface{})
	data["deviceId"] = device.GetID()
	adapter.send(AdapterRemoveDeviceRequest, data)

}

func (adapter *Adapter) cancelRemoveThing(deviceId string) {
	log.Info(fmt.Sprintf("adapter: %s execute pairing", adapter.id))
	data := make(map[string]interface{})
	data["deviceId"] = deviceId
	adapter.send(AdapterCancelRemoveDeviceCommand, data)
}

func (adapter *Adapter) getManager() *AddonManager {
	return adapter.plugin.pluginServer.manager
}

func (adapter *Adapter) send(messageType int, data map[string]interface{}) {
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
