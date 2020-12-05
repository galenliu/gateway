package addons

import (
	"fmt"
	messages "gitee.com/liu_guilin/WebThings-schema"
	"sync"
)

type AdapterProxy struct {
	*Adapter
	plugin  *Plugin
	manager *AddonsManager
	looker  *sync.Mutex
	name    string

	manifest interface{}
}

func NewAdapterProxy(manager *AddonsManager, adapterId string, name string, packetName string) *AdapterProxy {
	proxy := &AdapterProxy{}
	proxy.manager = manager
	//proxy.userProfile = manager.iGateway.GetUserProfile()
	//proxy.preferences = manager.iGateway.GetPreferences()
	proxy.ID = adapterId
	proxy.PackageName = packetName
	proxy.name = name
	return proxy
}

func (adapter *AdapterProxy) handlerDeviceAdded(dev *DeviceProxy) {
	adapter.looker.Lock()
	defer adapter.looker.Unlock()
	if dev.GetId() != "" {
		adapter.devices[dev.GetId()] = dev
	}
}

func (adapter *AdapterProxy) getDevice(devId string) *DeviceProxy {
	adapter.looker.Lock()
	defer adapter.looker.Unlock()
	device := adapter.devices[devId]
	return device
}

func (adapter *AdapterProxy) removeThing(dev *DeviceProxy) {
	log.Info(fmt.Sprintf("adapter delete thing Id: %v", dev.ID))
	adapter.sendMessage(messages.MessageTypeAdapterRemoveDeviceRequest, messages.AdapterRemoveDeviceRequest{
		AdapterId: adapter.ID,
		PluginId:  adapter.PackageName,
		DeviceId:  dev.ID,
	})
}

func (adapter *AdapterProxy) pairing() {

}

func (adapter *AdapterProxy) sendMessage(messageType int, data interface{}) {
	var message messages.BaseMessage
	message.MessageType = messageType
	message.Data = data
	adapter.plugin.sendMessage(message)
}
