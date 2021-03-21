package plugin

import (
	"addon"
	"context"
	"fmt"
	"gateway/pkg/log"
	"sync"
)

type pairingFunc func(ctx context.Context, cancelFunc func())

type AdapterProxy struct {
	*addon.Adapter
	pluginId       string
	plugin         *Plugin
	looker         *sync.Mutex
	isPairing      bool
	onPairingFunc  pairingFunc
	pairingContext context.Context
	manifest       interface{}
}

func NewAdapterProxy(manager *AddonManager, name, adapterId, pluginId, packageName string) *AdapterProxy {
	proxy := &AdapterProxy{}
	proxy.Adapter = addon.NewAdapter(manager, adapterId, name, packageName)
	proxy.pluginId = pluginId
	proxy.looker = new(sync.Mutex)

	return proxy
}

func (adapter *AdapterProxy) PropertyChanged(property, new *addon.Property) {
	property.Update(new)
}

func (adapter *AdapterProxy) handleSetPropertyValue(property *addon.Property, newValue interface{}) {
	data := make(map[string]interface{})
	data["deviceId"] = property.DeviceId
	data["propertyName"] = property.Name
	data["propertyValue"] = newValue
	adapter.send(AdapterCancelPairingCommand, data)
}

func (adapter *AdapterProxy) pairing(timeout float64) {
	log.Info(fmt.Sprintf("adapter: %s start pairing", adapter.ID))
	data := make(map[string]interface{})
	data["timeout"] = timeout
	adapter.send(AdapterStartPairingCommand, data)
}

func (adapter *AdapterProxy) cancelPairing() {
	log.Info(fmt.Sprintf("adapter: %s execute pairing", adapter.ID))
	data := make(map[string]interface{})
	adapter.send(AdapterCancelPairingCommand, data)
}

func (adapter *AdapterProxy) removeThing(device *addon.Device) {
	log.Info(fmt.Sprintf("adapter delete thing ID: %v", device.ID))
	data := make(map[string]interface{})
	data["deviceId"] = device.ID
	adapter.send(AdapterRemoveDeviceRequest, data)

}

func (adapter *AdapterProxy) cancelRemoveThing(deviceId string) {
	log.Info(fmt.Sprintf("adapter: %s execute pairing", adapter.ID))
	data := make(map[string]interface{})
	data["deviceId"] = deviceId
	adapter.send(AdapterCancelRemoveDeviceCommand, data)
}

func (adapter *AdapterProxy) getManager() *AddonManager {
	return adapter.plugin.pluginServer.manager
}

func (adapter *AdapterProxy) send(messageType int, data map[string]interface{}) {
	data["adapterId"] = adapter.ID
	adapter.plugin.send(messageType, data)
}
