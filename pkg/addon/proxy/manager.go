package proxy

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/manager"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	json "github.com/json-iterator/go"
	"sync"
	"time"
)

type Manager struct {
	*manager.Manager
	ipcClient   *IpcClient
	pluginId    string
	verbose     bool
	running     bool
	registered  bool
	userProfile *messages.PluginRegisterResponseJsonDataUserProfile
	preferences *messages.PluginRegisterResponseJsonDataPreferences
}

var once sync.Once
var instance *Manager

func NewAddonManager(pluginId string) (*Manager, error) {
	once.Do(
		func() {
			instance = &Manager{}
			instance.Manager = manager.NewManager()
			instance.pluginId = pluginId
			instance.verbose = true
			instance.registered = false
			instance.ipcClient = NewClient(instance, "9500")
			if instance.ipcClient != nil {
				instance.running = true
				instance.register()
			}
		},
	)
	if !instance.registered {
		return nil, fmt.Errorf("plugin not registered: %v", pluginId)
	}
	return instance, nil
}

func (m *Manager) RegisteredAdapter(adapters ...AdapterProxy) {
	for _, adapter := range adapters {
		adapter.Registered(m)
		m.AddAdapter(adapter)
		m.Send(messages.MessageType_AdapterAddedNotification, messages.AdapterAddedNotificationJsonData{
			AdapterId:   adapter.GetId(),
			Name:        adapter.GetName(),
			PackageName: adapter.GetPackageName(),
			PluginId:    m.pluginId,
		})
	}
}

func (m *Manager) RegisteredIntegration(integrations ...IntegrationProxy) {
	for _, ig := range integrations {

		m.AddIntegration(ig)
		m.Send(messages.MessageType_AdapterAddedNotification, messages.AdapterAddedNotificationJsonData{
			AdapterId:   ig.GetId(),
			Name:        ig.GetName(),
			PackageName: ig.GetPackageName(),
			PluginId:    m.pluginId,
		})
	}
}

func (m *Manager) HandleDeviceAdded(device DeviceProxy) {
	if m.verbose {
		fmt.Printf("manager device_added: %s \t\n", device.GetId())
	}
	m.Send(messages.MessageType_DeviceAddedNotification, messages.DeviceAddedNotificationJsonData{
		AdapterId: device.GetAdapter().GetId(),
		Device:    *device.ToMessage(),
		PluginId:  m.pluginId,
	})
}

func (m *Manager) HandleDeviceRemoved(device DeviceProxy) {
	if m.verbose {
		fmt.Printf("addon manager handle devices added, deviceId:%v\n", device.GetId())
	}
	m.Send(messages.MessageType_AdapterRemoveDeviceResponse, messages.AdapterRemoveDeviceResponseJsonData{
		AdapterId: device.GetAdapter().GetId(),
		DeviceId:  device.GetId(),
		PluginId:  m.pluginId,
	})
}

func (m *Manager) HandleAdapterRemoved(id string) {
	m.Manager.RemoveAdapter(id)
}

func (m *Manager) getAdapter(adapterId string) AdapterProxy {
	adapter := m.Manager.GetAdapter(adapterId)
	if adapter != nil {
		adp, ok := adapter.(AdapterProxy)
		if ok {
			return adp
		}
	}
	return nil
}

func (m *Manager) OnMessage(data []byte) {

	mt := json.Get(data, "messageType")
	dataNode := json.Get(data, "data")
	if mt.LastError() != nil || dataNode.LastError() != nil {
		fmt.Printf("message unmarshal err: %s", data)
		return
	}
	messageType := messages.MessageType(mt.ToInt())

	switch messageType {
	case messages.MessageType_PluginRegisterResponse:
		var msg messages.PluginRegisterResponseJsonData
		dataNode.ToVal(&msg)
		if dataNode.LastError() != nil {
			fmt.Printf("message unmarshal err:%s", dataNode.LastError().Error())
			return
		}
		m.registered = true
		m.preferences = &msg.Preferences
		m.userProfile = &msg.UserProfile
		return
	}
	if !m.registered {
		fmt.Printf("addon manager not registered")
		return
	}
	switch messageType {
	case messages.MessageType_PluginUnloadRequest:
		var msg messages.PluginUnloadRequestJsonData
		dataNode.ToVal(&msg)
		if dataNode.LastError() != nil {
			fmt.Printf("message unmarshal err:%s", dataNode.LastError().Error())
			return
		}
		m.Send(messages.MessageType_PluginUnloadResponse, messages.PluginUnloadResponseJsonData{PluginId: msg.PluginId})
		m.running = false
		go func() {
			time.AfterFunc(500*time.Millisecond, func() { m.Close() })
		}()
		return
	}

	adapterId := dataNode.Get("adapterId").ToString()
	adapter := m.getAdapter(adapterId)
	if adapter == nil {
		fmt.Printf("can not found adapter(%s)", adapterId)
		return
	}
	switch messageType {
	case messages.MessageType_AdapterStartPairingCommand:
		var msg messages.AdapterStartPairingCommandJsonData
		dataNode.ToVal(&msg)
		if dataNode.LastError() != nil {
			fmt.Printf("message unmarshal err:%s", dataNode.LastError().Error())
			return
		}
		go adapter.StartPairing(time.After(time.Duration(msg.Timeout) * time.Millisecond))
		return

	case messages.MessageType_AdapterCancelPairingCommand:
		var msg messages.AdapterCancelPairingCommandJsonData
		dataNode.ToVal(&msg)
		if dataNode.LastError() != nil {
			fmt.Printf("message unmarshal err:%s", dataNode.LastError().Error())
			return
		}
		go adapter.CancelPairing()
		return

	case messages.MessageType_AdapterUnloadRequest:
		var msg messages.AdapterUnloadRequestJsonData
		dataNode.ToVal(&msg)
		if dataNode.LastError() != nil {
			fmt.Printf("message unmarshal err:%s", dataNode.LastError().Error())
			return
		}
		go adapter.Unload()
		unloadFunc := func() {
			m.Send(messages.MessageType_AdapterUnloadResponse, messages.AdapterUnloadResponseJsonData{
				AdapterId: adapter.GetId(),
				PluginId:  m.pluginId,
			})
			m.HandleAdapterRemoved(adapter.GetId())
		}
		go unloadFunc()
		break
	}
	deviceId := dataNode.Get("deviceId").ToString()
	device := adapter.GetDevice(deviceId)
	if device == nil {
		fmt.Printf("device %s not found", deviceId)
		return
	}
	switch messageType {
	case messages.MessageType_AdapterCancelRemoveDeviceCommand:
		var msg messages.AdapterCancelRemoveDeviceCommandJsonData
		dataNode.ToVal(&msg)
		if dataNode.LastError() != nil {
			fmt.Printf("message unmarshal err:%s", dataNode.LastError().Error())
			return
		}
		go adapter.CancelRemoveThing(msg.DeviceId)
		return

	case messages.MessageType_DeviceSavedNotification:
		var msg messages.DeviceSavedNotificationJsonData
		dataNode.ToVal(&msg)
		if dataNode.LastError() != nil {
			fmt.Printf("message unmarshal err:%s", dataNode.LastError().Error())
			return
		}
		go adapter.HandleDeviceSaved(device)
		return

	case messages.MessageType_AdapterRemoveDeviceRequest:
		var msg messages.AdapterRemoveDeviceRequestJsonData
		dataNode.ToVal(&msg)
		if dataNode.LastError() != nil {
			fmt.Printf("message unmarshal err:%s", dataNode.LastError().Error())
			return
		}
		go adapter.HandleDeviceRemoved(device)
		return

	case messages.MessageType_DeviceSetPropertyCommand:
		var msg messages.DeviceSetPropertyCommandJsonData
		dataNode.ToVal(&msg)
		if dataNode.LastError() != nil {
			fmt.Printf("message unmarshal err:%s", dataNode.LastError().Error())
			return
		}
		prop := device.GetProperty(msg.PropertyName)
		if prop == nil {
			fmt.Printf("can not found propertyName(%s)", msg.PropertyName)
			return
		}
		go prop.SetValue(msg.PropertyValue)
		return

	case messages.MessageType_DeviceSetPinRequest:
		var msg messages.DeviceSetPinRequestJsonData
		dataNode.ToVal(&msg)
		if dataNode.LastError() != nil {
			fmt.Printf("message unmarshal err:%s", dataNode.LastError().Error())
			return
		}
		handleFunc := func() {
			err := device.SetPin(msg.Pin)
			var success = true
			if err != nil {
				fmt.Printf("device set pin err:%s", err.Error())
				success = false
			}
			m.Send(messages.MessageType_DeviceSetPinResponse, messages.DeviceSetPinResponseJsonData{
				AdapterId: device.GetAdapter().GetId(),
				Device:    device.ToMessage(),
				DeviceId:  &deviceId,
				MessageId: msg.MessageId,
				PluginId:  adapter.GetPackageName(),
				Success:   success,
			})
		}
		go handleFunc()

	case messages.MessageType_DeviceSetCredentialsRequest:
		var msg messages.DeviceSetCredentialsRequestJsonData
		dataNode.ToVal(&msg)
		if dataNode.LastError() != nil {
			fmt.Printf("message unmarshal err:%s", dataNode.LastError().Error())
			return
		}
		handleFunc := func() {
			err := device.SetCredentials(msg.Username, msg.Password)
			success := true
			if err != nil {
				success = false
				fmt.Printf(err.Error())
			}
			m.Send(messages.MessageType_DeviceSetCredentialsResponse, messages.DeviceSetCredentialsResponseJsonData{
				AdapterId: adapterId,
				Device:    nil,
				DeviceId:  &deviceId,
				MessageId: msg.MessageId,
				PluginId:  m.pluginId,
				Success:   success,
			})
		}
		go handleFunc()
		break
	}
}

func (m *Manager) sendConnectedStateNotification(device DeviceProxy, connected bool) {
	m.Send(messages.MessageType_DeviceConnectedStateNotification, messages.DeviceConnectedStateNotificationJsonData{
		AdapterId: device.GetAdapter().GetId(),
		Connected: connected,
		DeviceId:  device.GetId(),
		PluginId:  m.GetPluginId(),
	})
}

func (m *Manager) Send(messageType messages.MessageType, data any) {
	var message = struct {
		MessageType messages.MessageType `json:"messageType"`
		Data        any                  `json:"data"`
	}{MessageType: messageType, Data: data}
	m.ipcClient.Send(message)
}

func (m *Manager) register() {
	if !m.running {
		fmt.Printf("addon manager not running")
		return
	}
	m.Send(messages.MessageType_PluginRegisterRequest, messages.PluginRegisterRequestJsonData{PluginId: m.pluginId})
	time.Sleep(time.Duration(10) * time.Millisecond)
}

func (m *Manager) Close() {
	m.ipcClient.close()
	m.running = false
}

func (m *Manager) IsRunning() bool {
	if len(m.GetAdapters()) > 0 && m.running {
		return true
	}
	return false
}

func (m *Manager) GetPluginId() string {
	return m.pluginId
}

func (m *Manager) GetUserProfile() *messages.PluginRegisterResponseJsonDataUserProfile {
	return m.userProfile
}
