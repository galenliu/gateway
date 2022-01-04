package addon

import (
	"fmt"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	json "github.com/json-iterator/go"
	"log"
	"sync"
	"time"
)

type AdapterProxy interface {
	GetId() string
	GetName() string
	GetDevice(deviceId string) DeviceProxy
	Unload()
	CancelPairing()
	StartPairing(timeout time.Duration)
	HandleDeviceSaved(DeviceProxy)
	HandleDeviceRemoved(DeviceProxy)
}

type Manager struct {
	ipcClient   *IpcClient
	adapters    sync.Map
	packageName string
	verbose     bool
	running     bool
	registered  bool
	userProfile messages.PluginRegisterResponseJsonDataUserProfile
	preferences messages.PluginRegisterResponseJsonDataPreferences
}

var once sync.Once
var instance *Manager

func NewAddonManager(packageName string) (*Manager, error) {
	once.Do(
		func() {
			instance = &Manager{}
			instance.packageName = packageName
			instance.verbose = true
			instance.registered = false
			instance.ipcClient = NewClient(packageName, instance)
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
		m.adapters.Store(adapter.GetId(), adapter)
		m.send(messages.MessageType_AdapterAddedNotification, messages.AdapterAddedNotificationJsonData{
			AdapterId:   adapter.GetId(),
			Name:        adapter.GetName(),
			PackageName: m.packageName,
			PluginId:    m.packageName,
		})
	}
}

func (m *Manager) handleDeviceAdded(device DeviceProxy) {
	if m.verbose {
		log.Printf("addonManager: handle_device_added: %s", device.GetId())
	}
	m.send(messages.MessageType_DeviceAddedNotification, messages.DeviceAddedNotificationJsonData{
		AdapterId: device.GetAdapter().GetId(),
		Device:    messages.Device{},
		PluginId:  m.packageName,
	})
}

func (m *Manager) handleDeviceRemoved(device DeviceProxy) {
	if m.verbose {
		log.Printf("addon manager handle devices added, deviceId:%v\n", device.GetId())
	}
	m.send(messages.MessageType_AdapterRemoveDeviceResponse, messages.AdapterRemoveDeviceResponseJsonData{
		AdapterId: device.GetAdapter().GetId(),
		DeviceId:  device.GetId(),
		PluginId:  m.packageName,
	})
}

func (m *Manager) getAdapter(adapterId string) AdapterProxy {
	adapter, ok := m.adapters.Load(adapterId)
	if !ok {
		return nil
	}
	return adapter.(AdapterProxy)
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
		m.send(messages.MessageType_PluginUnloadResponse, messages.PluginUnloadResponseJsonData{PluginId: m.packageName})
		m.running = false
		var closeFun = func() {
			time.AfterFunc(500*time.Millisecond, func() { m.close() })
		}
		go closeFun()
		return
	}

	var adapterId = json.Get(data, "data", "adapterId").ToString()
	adapter := m.getAdapter(adapterId)
	if adapter == nil {
		log.Printf("can not found adapter(%s)", adapterId)
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
			proxy.send(messages.MessageType_AdapterUnloadResponse, messages.AdapterUnloadResponseJsonData{
				AdapterId: adapter.GetId(),
				PluginId:  m.packageName,
			})
		}
		go unloadFunc(m, adapter)
		m.adapters.Delete(adapter.GetId())
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
		log.Printf(adapter.GetId())

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
			log.Printf("can not found propertyName(%s)", propName)
			return
		}
		propChanged := func(newValue any) error {
			prop.SetValue(newValue)
			return nil
		}
		e := propChanged(newValue)
		if e != nil {
			log.Printf(e.Error())
			return
		}
		p := prop.ToDescription()
		m.send(messages.MessageType_DevicePropertyChangedNotification, messages.DevicePropertyChangedNotificationJsonData{
			AdapterId: adapter.GetId(),
			DeviceId:  device.GetId(),
			PluginId:  m.packageName,
			Property: messages.Property{
				Type:        p.Type,
				AtType:      p.AtType,
				Description: p.Description,
				Enum:        p.Enum,
				Maximum:     p.Maximum,
				Minimum:     p.Minimum,
				MultipleOf:  p.MultipleOf,
				Name:        p.Name,
				ReadOnly:    p.ReadOnly,
				Title:       p.Title,
				Unit:        p.Unit,
				Value:       p.Value,
			},
		})

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
				m.send(messages.MessageType_DeviceSetCredentialsResponse, messages.DeviceSetCredentialsResponseJsonData{
					AdapterId: adapter.GetId(),
					Device:    nil,
					DeviceId: func(s string) *string {
						if s == "" {
							return nil
						}
						return &s
					}(device.GetId()),
					MessageId: messageId,
					PluginId:  m.packageName,
					Success:   true,
				})
				fmt.Printf(err.Error())
				return
			}

			m.send(messages.MessageType_DeviceSetCredentialsResponse, messages.DeviceSetCredentialsResponseJsonData{
				AdapterId: adapter.GetId(),
				Device:    nil,
				DeviceId: func(s string) *string {
					if s == "" {
						return nil
					}
					return &s
				}(device.GetId()),
				MessageId: messageId,
				PluginId:  m.packageName,
				Success:   false,
			})
			return
		}
		go handleFunc()
		break
	}
}

func (m *Manager) sendConnectedStateNotification(device DeviceProxy, connected bool) {
	m.send(messages.MessageType_DeviceConnectedStateNotification, messages.DeviceConnectedStateNotificationJsonData{
		AdapterId: device.GetAdapter().GetId(),
		Connected: connected,
		DeviceId:  device.GetId(),
		PluginId:  m.packageName,
	})
}

func (m *Manager) send(messageType messages.MessageType, data any) {
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
	m.send(messages.MessageType_PluginRegisterRequest, messages.PluginRegisterRequestJsonData{PluginId: m.packageName})
}

func (m *Manager) close() {
	m.ipcClient.close()
	m.running = false
}
