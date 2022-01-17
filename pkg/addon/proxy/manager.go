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
	userProfile messages.PluginRegisterResponseJsonDataUserProfile
	preferences messages.PluginRegisterResponseJsonDataPreferences
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
			instance.ipcClient = NewClient(pluginId, instance)
			if instance.ipcClient != nil {
				instance.running = true
			}
			instance.register()
		},
	)
	if instance.ipcClient == nil {
		return nil, fmt.Errorf("ipc client not available")
	}
	instance.running = true
	return instance, nil
}

func (m *Manager) AddAdapters(adapters ...AdapterProxy) {
	for _, adapter := range adapters {
		adapter.SetManager(m)
		m.AddAdapter(adapter)
		m.Send(messages.MessageType_AdapterAddedNotification, messages.AdapterAddedNotificationJsonData{
			AdapterId:   adapter.GetId(),
			Name:        adapter.GetName(),
			PackageName: adapter.GetPackageName(),
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
		Device:    device.ToMessage(),
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

func (m *Manager) onMessage(data []byte) {

	var messageType = messages.MessageType(json.Get(data, "messageType").ToInt())
	var dataAny = json.Get(data, "data")

	switch messageType {
	case messages.MessageType_PluginRegisterResponse:
		var msg messages.PluginRegisterResponseJsonData
		dataAny.ToVal(&msg)
		if &m != nil {
			m.registered = true
			m.preferences = msg.Preferences
			m.userProfile = msg.UserProfile
		}
		return
	}
	if !m.registered {
		fmt.Printf("addon manager not registered")
		return
	}
	switch messageType {

	case messages.MessageType_PluginUnloadRequest:
		m.Send(messages.MessageType_PluginUnloadResponse, messages.PluginUnloadResponseJsonData{PluginId: m.pluginId})
		m.running = false
		var closeFun = func() {
			time.AfterFunc(500*time.Millisecond, func() { m.Close() })
		}
		go closeFun()
		return
	}

	var adapterId = json.Get(data, "data", "adapterId").ToString()
	adapter := m.getAdapter(adapterId)
	if adapter == nil {
		fmt.Printf("can not found adapter(%s)", adapterId)
		return
	}

	switch messageType {
	//adapter pairing command
	case messages.MessageType_AdapterStartPairingCommand:
		timeout := json.Get(data, "data", "timeout").ToFloat64()
		go adapter.StartPairing(time.Duration(timeout) * time.Millisecond)
		return

	case messages.MessageType_AdapterCancelPairingCommand:
		go adapter.CancelPairing()
		return

		//adapter unload request

	case messages.MessageType_AdapterUnloadRequest:
		go adapter.Unload()
		unloadFunc := func(proxy *Manager, adapter AdapterProxy) {
			data := make(map[string]any)
			data["adapterId"] = adapter.GetId()
			proxy.Send(messages.MessageType_AdapterUnloadResponse, messages.AdapterUnloadResponseJsonData{
				AdapterId: adapter.GetId(),
				PluginId:  m.pluginId,
			})
		}
		go unloadFunc(m, adapter)
		m.RemoveAdapter(adapter.GetId())
		break
	}
	var deviceId = json.Get(data, "data", "deviceId").ToString()
	device := adapter.GetDevice(deviceId)
	if device == nil {
		return
	}

	switch messageType {
	case messages.MessageType_AdapterCancelRemoveDeviceCommand:
		adapter := m.getAdapter(adapterId)
		fmt.Printf(adapter.GetId())

	case messages.MessageType_DeviceSavedNotification:

		go adapter.HandleDeviceSaved(device)
		return

		//adapter remove devices request

	case messages.MessageType_AdapterRemoveDeviceRequest:
		adapter.HandleDeviceRemoved(device)

		//devices set properties command

	case messages.MessageType_DeviceSetPropertyCommand:
		propName := json.Get(data, "data", "propertyName").ToString()
		newValue := json.Get(data, "data", "propertyValue").GetInterface()
		prop := device.GetProperty(propName)
		if prop == nil {
			fmt.Printf("can not found propertyName(%s)", propName)
			return
		}
		propChanged := func(newValue any) error {
			prop.SetValue(newValue)
			return nil
		}
		e := propChanged(newValue)
		if e != nil {
			fmt.Printf(e.Error())
			return
		}
		return

	case messages.MessageType_DeviceSetPinRequest:
		//var pin PIN
		//pin.Pattern = json.Get(data, "data", "pin", "pattern").GetInterface()
		//pin.Required = json.Get(data, "data", "pin", "required").ToBool()
		//messageId := json.Get(data, "data", "message_id").ToInt()
		//if messageId == 0 {
		//	log.Fatal("DeviceSetPinRequest:  non  messageId")
		//}
		//
		//handleFunc := func() {
		//	data := make(map[string]interface{})
		//	data[Aid] = adapterId
		//	data[Did] = deviceId
		//	data["devx"] = device
		//	data["messageId"] = messageId
		//	err := device.SetPin(pin)
		//	if err == nil {
		//		data["success"] = true
		//		m.send(DeviceSetPinResponse, data)
		//
		//	} else {
		//		data["success"] = false
		//		m.send(DeviceSetPinResponse, data)
		//	}
		//}
		//go handleFunc()

	case messages.MessageType_DeviceSetCredentialsRequest:
		messageId := json.Get(data, "data", "messageId").ToInt()
		username := json.Get(data, "data", "username").ToString()
		password := json.Get(data, "data", "password").ToString()

		handleFunc := func() {
			err := device.SetCredentials(username, password)
			if err != nil {
				m.Send(messages.MessageType_DeviceSetCredentialsResponse, messages.DeviceSetCredentialsResponseJsonData{
					AdapterId: adapter.GetId(),
					Device:    nil,
					DeviceId: func(s string) *string {
						if s == "" {
							return nil
						}
						return &s
					}(device.GetId()),
					MessageId: messageId,
					PluginId:  m.pluginId,
					Success:   true,
				})
				fmt.Printf(err.Error())
				return
			}

			m.Send(messages.MessageType_DeviceSetCredentialsResponse, messages.DeviceSetCredentialsResponseJsonData{
				AdapterId: adapter.GetId(),
				Device:    nil,
				DeviceId: func(s string) *string {
					if s == "" {
						return nil
					}
					return &s
				}(device.GetId()),
				MessageId: messageId,
				PluginId:  m.pluginId,
				Success:   false,
			})
			return
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
		PluginId:  m.pluginId,
	})
}

func (m *Manager) Send(messageType messages.MessageType, data any) {
	var message = struct {
		MessageType messages.MessageType `json:"messageType"`
		Data        any                  `json:"data"`
	}{MessageType: messageType, Data: data}
	m.ipcClient.send(message)
}

func (m *Manager) register() {
	if !m.running {
		fmt.Printf("addon manager not running")
		return
	}
	m.Send(messages.MessageType_PluginRegisterRequest, messages.PluginRegisterRequestJsonData{PluginId: m.pluginId})
}

func (m *Manager) Close() {
	m.ipcClient.close()
	m.running = false
}

func (m *Manager) IsRunning() bool {
	return m.running
}

func (m *Manager) GetPluginId() string {
	return m.pluginId
}
