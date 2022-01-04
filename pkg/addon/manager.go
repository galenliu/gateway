package addon

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/devices"
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
}

var once sync.Once
var instance *Manager

func NewAddonManager(packageName string) *Manager {
	once.Do(
		func() {
			instance = &Manager{}
			instance.packageName = packageName
			instance.running = true
			instance.verbose = true
			instance.ipcClient = NewClient(packageName, instance)
		},
	)
	return instance
}

func (m *Manager) AddAdapters(adapters ...AdapterProxy) {
	if m.running {
		m.run()
	}
	for _, adapter := range adapters {
		m.adapters.Store(adapter.GetId(), adapter)
		adapter.GetId()
		data := make(map[string]any)
		data["adapterId"] = adapter.GetId()
		data["name"] = adapter.GetName
		data["packageName"] = m.packageName
		m.send(AdapterAddedNotification, data)
	}
}

func (m *Manager) handleDeviceAdded(device DeviceProxy) {
	if m.verbose {
		log.Printf("addonManager: handle_device_added: %s", device.GetId())
	}
	data := make(map[string]any)
	data["adapterId"] = device.GetId()
	description, err := json.Marshal(device)
	if err != nil {
		return
	}
	data["device"] = description
	m.send(DeviceAddedNotification, data)
}

func (m *Manager) handleDeviceRemoved(device DeviceProxy) {
	if m.verbose {
		log.Printf("addon manager handle devices added, deviceId:%v\n", device.GetId())
	}
	data := make(map[string]any)
	data["adapterId"] = device.GetAdapter()
	data["deviceId"] = device.GetId()

	m.send(AdapterRemoveDeviceResponse, data)
}

func (m *Manager) getAdapter(adapterId string) AdapterProxy {
	adapter, ok := m.adapters.Load(adapterId)
	if !ok {
		return nil
	}
	return adapter.(AdapterProxy)
}

func (m *Manager) onMessage(data []byte) {

	var messageType = json.Get(data, "messageType").ToInt()

	switch messageType {
	//卸载plugin
	case PluginUnloadRequest:
		data := make(map[string]any)
		m.send(PluginUnloadResponse, data)
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
	case AdapterStartPairingCommand:
		timeout := json.Get(data, "data", "timeout").ToFloat64()
		go adapter.StartPairing(time.Duration(timeout) * time.Millisecond)
		return

	case AdapterCancelPairingCommand:
		go adapter.CancelPairing()
		return

		//adapter unload request

	case AdapterUnloadRequest:
		go adapter.Unload()
		unloadFunc := func(proxy *Manager, adapter AdapterProxy) {
			data := make(map[string]any)
			data["adapterId"] = adapter.GetId()
			proxy.send(AdapterUnloadResponse, data)
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
	case AdapterCancelRemoveDeviceCommand:
		adapter := m.getAdapter(adapterId)
		log.Printf(adapter.GetId())

	case DeviceSavedNotification:

		go adapter.HandleDeviceSaved(device)
		return

		//adapter remove devices request

	case AdapterRemoveDeviceRequest:
		adapter.HandleDeviceRemoved(device)

		//devices set properties command

	case DeviceSetPropertyCommand:
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
		data := make(map[string]any)
		data["adapterId"] = adapterId
		data["deviceId"] = device.GetId()
		//data["property"] = prop.AsDict()
		m.send(DevicePropertyChangedNotification, data)

	case DeviceSetPinRequest:
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

	case DeviceSetCredentialsRequest:
		messageId := json.Get(data, "data", "messageId").ToInt()
		username := json.Get(data, "data", "username").ToString()
		password := json.Get(data, "data", "password").ToString()

		handleFunc := func() {
			err := device.SetCredentials(username, password)
			data := make(map[string]any)
			data["adapterId"] = adapterId
			data["deviceId"] = deviceId
			data["messageId"] = messageId
			if err != nil {
				data["success"] = true
				m.send(DeviceSetCredentialsResponse, data)
				fmt.Printf(err.Error())
				return
			}
			data["success"] = false
			m.send(DeviceSetCredentialsResponse, data)
			return
		}
		go handleFunc()
		break
	}
}

func (m *Manager) sendConnectedStateNotification(device *devices.Device, connected bool) {
	data := make(map[string]any)
	data["adapterId"] = ""
	data["deviceId"] = device.GetId()
	data["connected"] = connected
	m.send(DeviceConnectedStateNotification, data)
}

func (m *Manager) run() {
	m.ipcClient.Register()
}

func (m *Manager) send(messageType messages.MessageType, data any) {
	var message = struct {
		MessageType messages.MessageType `json:"messageType"`
		Data        any                  `json:"data"`
	}{MessageType: messageType, Data: data}
	d, er := json.MarshalIndent(message, "", " ")
	if er != nil {
		log.Fatal(er)
		return
	}
	m.ipcClient.sendMessage(d)
}

func (m *Manager) register() {
	m.ipcClient.sendMessage()
}

func (m *Manager) close() {
	m.ipcClient.close()
	m.running = false
}
