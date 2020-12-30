package addons

import (
	"fmt"
	"gateway/pkg/log"
	addon "gitee.com/liu_guilin/gateway-addon-golang"
	json "github.com/json-iterator/go"
	"sync"
)

type AdapterProxy struct {
	*addon.Adapter
	plugin       *Plugin
	addonManager *AddonManager
	looker       *sync.Mutex
	manifest     interface{}
}

func NewAdapterProxy(manager *AddonManager, plugin *Plugin, adapterId string, name, packageName string) *AdapterProxy {
	proxy := &AdapterProxy{}
	proxy.PackageName = packageName
	proxy.plugin = plugin
	proxy.addonManager = manager
	proxy.ID = adapterId
	proxy.PackageName = name
	proxy.looker = new(sync.Mutex)
	return proxy
}

func (adapter *AdapterProxy) handlerDeviceAdded(dev *DeviceProxy) {
	adapter.addonManager.handlerDeviceAdded(dev)

}

func (adapter *AdapterProxy) removeThing(dev *DeviceProxy) {
	log.Info(fmt.Sprintf("adapter delete thing Id: %v", dev.ID))
	adapter.sendMessage(AdapterRemoveDeviceRequest, struct {
		AdapterId string
		PluginId  string
		DeviceId  string
	}{
		AdapterId: adapter.ID,
		PluginId:  adapter.PackageName,
		DeviceId:  dev.ID,
	})
}

func (adapter *AdapterProxy) Pairing(timeout int) {
	log.Info(fmt.Sprintf("adapter: %s start pairing", adapter.ID))

	adapter.sendMessage(AdapterStartPairingCommand, struct {
		PluginID  string `json:"pluginId"`
		AdapterID string `json:"adapterId"`
		Timeout   int    `json:"timeout"`
	}{
		PluginID:  adapter.PackageName,
		AdapterID: adapter.ID,
		Timeout:   timeout})
}

func (adapter *AdapterProxy) sendMessage(messageType int, msg interface{}) {

	message := struct {
		MessageType int         `json:"messageType"`
		Data        interface{} `json:"data"`
	}{
		MessageType: messageType,
		Data:        msg,
	}
	bt, err := json.Marshal(message)
	if err != nil {
		return
	}
	adapter.plugin.sendData(bt)

}
