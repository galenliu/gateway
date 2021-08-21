package plugin

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rpc"
	"github.com/galenliu/gateway/plugin/internal"
	"sync"
)

type pairingFunc func(ctx context.Context, cancelFunc func())

type managerProxy interface {
	handleDeviceAdded(device *internal.Device)
}

type Adapter struct {
	ID          string
	name        string
	pluginId    string
	plugin      *Plugin
	looker      *sync.Mutex
	isPairing   bool
	devices     sync.Map
	packageName string
	logger      logging.Logger
}

func NewAdapter(p *Plugin, name, adapterId, packageName string, log logging.Logger) *Adapter {
	adapter := &Adapter{}
	adapter.logger = log
	adapter.ID = adapterId
	adapter.pluginId = p.pluginId
	adapter.name = name
	adapter.packageName = packageName
	adapter.plugin = p
	adapter.looker = new(sync.Mutex)
	return adapter
}

func (adapter *Adapter) pairing(timeout float64) {
	adapter.logger.Info(fmt.Sprintf("adapter: %s start pairing", adapter.ID))
	data := make(map[string]interface{})
	data["timeout"] = timeout
	adapter.sendMessage(rpc.MessageType_AdapterStartPairingCommand, data)
}

func (adapter *Adapter) cancelPairing() {
	adapter.logger.Info(fmt.Sprintf("adapter: %s execute pairing", adapter.ID))
	data := make(map[string]interface{})
	adapter.sendMessage(rpc.MessageType_AdapterCancelPairingCommand, data)
}

func (adapter *Adapter) removeThing(device *internal.Device) {
	adapter.logger.Info("adapter delete thing Id: %v", device.ID)
	data := make(map[string]interface{})
	data["deviceId"] = device.ID
	adapter.sendMessage(internal.AdapterRemoveDeviceRequest, data)

}

func (adapter *Adapter) cancelRemoveThing(deviceId string) {
	adapter.logger.Info(fmt.Sprintf("adapter: %s execute pairing", adapter.ID))
	data := make(map[string]interface{})
	data["deviceId"] = deviceId
	adapter.sendMessage(rpc.MessageType_AdapterCancelRemoveDeviceCommand, data)
}

func (adapter *Adapter) sendMessage(messageType rpc.MessageType, data map[string]interface{}) {
	data["adapterId"] = adapter.ID
	adapter.plugin.SendMessage(messageType, data)
}

func (adapter *Adapter) handleDeviceRemoved(d *Device) {
	adapter.devices.Delete(d.ID)
	adapter.plugin.pluginServer.manager.handleDeviceRemoved(d)
}

func (adapter *Adapter) handleDeviceAdded(device *Device) {
	adapter.devices.Store(device.ID, device)
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
