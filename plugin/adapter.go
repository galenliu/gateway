package plugin

import (
	"context"
	"fmt"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"sync"
)

type Adapter struct {
	id                 string
	name               string
	isPairing          bool
	packageName        string
	logger             logging.Logger
	plugin             *Plugin
	devices            sync.Map
	setCredentialsTask sync.Map
	setPinTask         sync.Map
	nextId             int
}

func NewAdapter(plugin *Plugin, adapterId, name, packageName string, log logging.Logger) *Adapter {
	adapter := &Adapter{}
	adapter.plugin = plugin
	adapter.logger = log
	adapter.id = adapterId
	adapter.name = name
	adapter.packageName = packageName
	adapter.nextId = 0
	return adapter
}

func (adapter *Adapter) setDeviceCredentials(ctx context.Context, thingId, username, password string) (*messages.Device, error) {
	messageId := adapter.generatedId()
	t, ok := adapter.setCredentialsTask.LoadOrStore(messageId, make(chan *messages.Device))
	defer adapter.setCredentialsTask.Delete(messageId)
	adapter.send(messages.MessageType_DeviceSetCredentialsRequest,
		messages.DeviceSetCredentialsRequestJsonData{
			AdapterId: adapter.getId(),
			DeviceId:  thingId,
			MessageId: messageId,
			Password:  password,
			PluginId:  adapter.plugin.getId(),
			Username:  username,
		})
	task, ok := t.(chan *messages.Device)
	if ok {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout for setDeviceCredentials")
		case d := <-task:
			return d, nil
		}
	}
	return nil, fmt.Errorf("failed to set device credentials")
}

func (adapter *Adapter) setDevicePin(ctx context.Context, deviceId, pin string) (*messages.Device, error) {
	messageId := adapter.generatedId()
	t, ok := adapter.setPinTask.LoadOrStore(messageId, make(chan *messages.Device))
	defer adapter.setPinTask.Delete(messageId)
	adapter.send(messages.MessageType_DeviceSetPinRequest,
		messages.DeviceSetPinRequestJsonData{
			AdapterId: adapter.getId(),
			DeviceId:  deviceId,
			MessageId: 0,
			Pin:       pin,
			PluginId:  adapter.getPlugin().getId(),
		})
	task, ok := t.(chan *messages.Device)
	if ok {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout for setting device pin")
		case d := <-task:
			return d, nil
		}
	}
	return nil, fmt.Errorf("failed set device pin")
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
	adapter.logger.Info("adapter: %s id: %s CancelRemoveThing:", adapter.getName(), adapter.getId(), device.GetId())
	adapter.send(messages.MessageType_AdapterCancelRemoveDeviceCommand, data)
}

func (adapter *Adapter) send(messageType messages.MessageType, data interface{}) {
	adapter.plugin.send(messageType, data)
}

func (adapter *Adapter) handleDeviceRemoved(d *Device) {
	adapter.devices.Delete(d.GetId())
	adapter.plugin.manager.handleDeviceRemoved(d)
}

func (adapter *Adapter) handleDeviceAdded(device *Device) {
	adapter.devices.Store(device.GetId(), device)
	adapter.plugin.manager.handleDeviceAdded(device)
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
	for _, device := range adapter.getDevices() {
		adapter.handleDeviceRemoved(device)
	}
}

func (adapter *Adapter) generatedId() int {
	adapter.nextId = adapter.nextId + 1
	return adapter.nextId
}
