package plugin

import (
	"addon"
	"context"
	"fmt"
	"gateway/pkg/log"
	json "github.com/json-iterator/go"
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

func (adapter *AdapterProxy) setPropertyValue(property *addon.Property, newValue interface{}) {
	adapter.sendMessage(DeviceSetPropertyCommand, struct {
		AdapterId     string      `json:"adapterId"`
		PluginId      string      `json:"pluginId"`
		DeviceId      string      `json:"deviceId"`
		PropertyName  string      `json:"propertyName"`
		PropertyValue interface{} `json:"propertyValue"`
	}{
		AdapterId:     adapter.ID,
		PluginId:      adapter.pluginId,
		DeviceId:      property.DeviceId,
		PropertyName:  property.Name,
		PropertyValue: newValue,
	})
}

func (adapter *AdapterProxy) handleRemoveThing(device *addon.Device) {
	log.Info(fmt.Sprintf("adapter delete thing ID: %v", device.ID))
	adapter.sendMessage(AdapterRemoveDeviceRequest, struct {
		AdapterId string `json:"adapterId"`
		PluginId  string `json:"pluginId"`
		DeviceId  string `json:"deviceId"`
	}{
		AdapterId: adapter.ID,
		PluginId:  adapter.pluginId,
		DeviceId:  device.ID,
	})
}

func (adapter *AdapterProxy) pairing(timeout float64) {
	log.Info(fmt.Sprintf("adapter: %s start pairing", adapter.ID))
	adapter.sendMessage(AdapterStartPairingCommand, struct {
		PluginId  string  `json:"pluginId"`
		AdapterID string  `json:"adapterId"`
		Timeout   float64 `json:"timeout"`
	}{
		PluginId:  adapter.pluginId,
		AdapterID: adapter.ID,
		Timeout:   timeout})
}

func (adapter *AdapterProxy) cancelPairing() {
	log.Info(fmt.Sprintf("adapter: %s start pairing", adapter.ID))
	adapter.sendMessage(AdapterCancelPairingCommand, struct {
		PluginId  string `json:"pluginId"`
		AdapterID string `json:"adapterId"`
	}{
		PluginId:  adapter.pluginId,
		AdapterID: adapter.ID,
	})
}

func (adapter *AdapterProxy) getManager() *AddonManager {
	return adapter.plugin.pluginServer.manager
}

type tag struct {
	AdapterId string
	PluginId  string
}

func (adapter *AdapterProxy) sendMessage(messageType int, msg interface{}) {
	message := struct {
		MessageType int         `json:"messageType"`
		Data        interface{} `json:"data"`
	}{
		MessageType: messageType,
		Data:        msg,
	}
	bt, err := json.MarshalIndent(message, "", " ")
	if err != nil {
		return
	}
	adapter.plugin.sendData(bt)
}
