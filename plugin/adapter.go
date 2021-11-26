package plugin

import (
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"sync"
)

type Adapter struct {
	id          string
	name        string
	looker      *sync.Mutex
	isPairing   bool
	devices     sync.Map
	packageName string
	logger      logging.Logger
	plugin      *Plugin
}

func NewAdapter(plugin *Plugin, adapterId, name, packageName string, log logging.Logger) *Adapter {
	adapter := &Adapter{}
	adapter.plugin = plugin
	adapter.logger = log
	adapter.id = adapterId
	adapter.name = name
	adapter.packageName = packageName
	adapter.looker = new(sync.Mutex)
	return adapter
}

func (adapter *Adapter) startPairing(timeout int) {
	data := messages.AdapterStartPairingCommandJsonData{
		AdapterId: adapter.getId(),
		PluginId:  adapter.plugin.getId(),
		Timeout:   timeout,
	}
	adapter.logger.Infof("adapter %s startPairing", adapter.getId())
	adapter.send(messages.MessageType_AdapterStartPairingCommand, data)
}

func (adapter *Adapter) cancelPairing() {
	data := messages.AdapterCancelPairingCommandJsonData{
		AdapterId: adapter.getId(),
		PluginId:  adapter.plugin.getId(),
	}
	adapter.logger.Infof("adapter %s cancel startPairing", adapter.id)
	adapter.send(messages.MessageType_AdapterCancelPairingCommand, data)
}

func (adapter *Adapter) removeThing(device *Device) {
	data := messages.AdapterRemoveDeviceResponseJsonData{
		AdapterId: adapter.getId(),
		DeviceId:  device.GetId(),
		PluginId:  adapter.plugin.getId(),
	}
	adapter.logger.Infof("adapter delete thing Id: %v", device.GetId())
	adapter.send(messages.MessageType_AdapterRemoveDeviceRequest, data)
}

func (adapter *Adapter) cancelRemoveThing(device *Device) {
	data := messages.AdapterCancelRemoveDeviceCommandJsonData{
		AdapterId: adapter.getId(),
		DeviceId:  device.GetId(),
		PluginId:  adapter.plugin.getId(),
	}
	adapter.logger.Info("adapter: %s id: %s cancelRemoveThing:", adapter.getName(), adapter.getId(), device.GetId())
	adapter.send(messages.MessageType_AdapterCancelRemoveDeviceCommand, data)
}

func (adapter *Adapter) sendMsg(messageType messages.MessageType, data map[string]interface{}) {
	data["adapterId"] = adapter.id
	adapter.plugin.sendMsg(messageType, data)
}

func (adapter *Adapter) send(messageType messages.MessageType, data interface{}) {
	adapter.plugin.send(messageType, data)
}

func (adapter *Adapter) handleDeviceRemoved(d *Device) {
	adapter.devices.Delete(d.GetId())
	adapter.plugin.pluginServer.manager.handleDeviceRemoved(d)
}

func (adapter *Adapter) handleDeviceAdded(device *Device) {
	adapter.devices.Store(device.GetId(), device)
	adapter.plugin.pluginServer.manager.handleDeviceAdded(device)
}

func (adapter *Adapter) getDevice(deviceId string) *Device {
	d, ok := adapter.devices.Load(deviceId)
	device, ok := d.(*Device)
	if !ok {
		return nil
	}
	return device
}

func (adapter *Adapter) getDevices() (devices []*Device) {
	adapter.devices.Range(func(key, value interface{}) bool {
		device, ok := value.(*Device)
		if ok {
			devices = append(devices, device)
		}
		return true
	})
	return
}

func (adapter *Adapter) getName() string {
	return adapter.name
}

func (adapter *Adapter) getId() string {
	return adapter.id
}

func (adapter *Adapter) getPlugin() *Plugin {
	return adapter.plugin
}

func (adapter *Adapter) unload() {

}
