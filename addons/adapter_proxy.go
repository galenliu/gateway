package addons

import (
	"fmt"
	"gateway/pkg/log"
	json "github.com/json-iterator/go"
	"sync"
)

type AdapterProxy struct {
	ID          string
	PackageName string
	Name        string
	plugin      *Plugin
	manager     *AddonsManager
	looker      *sync.Mutex
	manifest    interface{}
}

func NewAdapterProxy(manager *AddonsManager, plugin *Plugin, adapterId string, name, packetName string) *AdapterProxy {
	proxy := &AdapterProxy{}

	proxy.manager = manager
	proxy.plugin = plugin
	proxy.ID = adapterId
	proxy.Name = name
	proxy.looker = new(sync.Mutex)
	proxy.PackageName = packetName
	return proxy
}

func (adapter *AdapterProxy) handlerDeviceAdded(dev *DeviceProxy) {
	adapter.manager.handlerDeviceAdded(dev)
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
