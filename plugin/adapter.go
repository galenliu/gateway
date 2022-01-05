package plugin

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/adapter"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"sync"
)

type Adapter struct {
	adapter.Adapter
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
	a := &Adapter{}
	a.Adapter = adapter.Adapter{
		Id:   adapterId,
		Name: name,
	}
	a.plugin = plugin
	a.logger = log
	a.Id = adapterId
	a.packageName = packageName
	a.nextId = 0
	return a
}

func (adapter *Adapter) setDeviceCredentials(ctx context.Context, thingId, username, password string) (*messages.Device, error) {
	messageId := adapter.generatedId()
	t, ok := adapter.setCredentialsTask.LoadOrStore(messageId, make(chan *messages.Device))
	defer adapter.setCredentialsTask.Delete(messageId)
	adapter.send(messages.MessageType_DeviceSetCredentialsRequest,
		messages.DeviceSetCredentialsRequestJsonData{
			AdapterId: adapter.GetId(),
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
			AdapterId: adapter.GetId(),
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
		AdapterId: adapter.GetId(),
		PluginId:  adapter.plugin.getId(),
		Timeout:   timeout,
	}
	adapter.logger.Infof("adapter %s startPairing", adapter.GetId())
	adapter.send(messages.MessageType_AdapterStartPairingCommand, data)
}

func (adapter *Adapter) cancelPairing() {
	data := messages.AdapterCancelPairingCommandJsonData{
		AdapterId: adapter.GetId(),
		PluginId:  adapter.plugin.getId(),
	}
	adapter.logger.Infof("adapter %s cancel startPairing", adapter.GetId())
	adapter.send(messages.MessageType_AdapterCancelPairingCommand, data)
}

func (adapter *Adapter) removeThing(device *Device) {
	data := messages.AdapterRemoveDeviceResponseJsonData{
		AdapterId: adapter.GetId(),
		DeviceId:  device.GetId(),
		PluginId:  adapter.plugin.getId(),
	}
	adapter.logger.Infof("adapter delete thing Id: %v", device.GetId())
	adapter.send(messages.MessageType_AdapterRemoveDeviceRequest, data)
}

func (adapter *Adapter) cancelRemoveThing(device *Device) {
	data := messages.AdapterCancelRemoveDeviceCommandJsonData{
		AdapterId: adapter.GetId(),
		DeviceId:  device.GetId(),
		PluginId:  adapter.plugin.getId(),
	}
	adapter.logger.Info("adapter: %s id: %s CancelRemoveThing:", adapter.GetName(), adapter.GetId(), device.GetId())
	adapter.send(messages.MessageType_AdapterCancelRemoveDeviceCommand, data)
}

func (adapter *Adapter) send(messageType messages.MessageType, data any) {
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
	adapter.devices.Range(func(key, value any) bool {
		device, ok := value.(*Device)
		if ok {
			devices = append(devices, device)
		}
		return true
	})
	return
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
