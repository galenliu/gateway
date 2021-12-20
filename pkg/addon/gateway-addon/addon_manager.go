package addon

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/gateway-addon/devices"
	json "github.com/json-iterator/go"
	"log"
	"sync"
	"time"
)

const (
	Aid = "adapterId"
	Pid = "pluginId"
	Did = "deviceId"
)

type IProperty interface {
	SetValue(any)
	GetValue() any
	GetName() string
	SetName(string)
	GetAtType() string
	GetType() string
	AsDict() []byte
	ToValue(any) any

	DoPropertyChanged(string)
	UpdateProperty(string)

	SetDeviceProxy(device IDevice)
}

type IAction interface {
	MarshalJson() []byte
}

type IEvent interface {
	MarshalJson() []byte
}

type IDevice interface {
	Send(int, map[string]any)
	GetProperty(name string) IProperty
	SetCredentials(username string, password string) error
	SetPin(pin any) error
	GetDescription() string
	ToJson() string
	GetID() string
	AsDict() Map
	GetAdapterId() string
}

type IAdapter interface {
	//SendPairingPrompt(promt, url string, device *Device)
	StartPairing(timeout float64)
	CancelPairing()
	GetID() string
	GetName() string
	Unload()
	HandleDeviceSaved(device IDevice)
	HandleDeviceRemoved(device IDevice)

	Send(mt int, data map[string]any)

	getDevice(deviceId string) IDevice
	setManager(manager *AddonManager)
}

type AddonManager struct {
	ipcClient   *IpcClient
	adapters    map[string]IAdapter ``
	packageName string
	verbose     bool
	running     bool
}

var once sync.Once
var instance *AddonManager

func InitAddonManager(packageName string) *AddonManager {
	once.Do(
		func() {
			instance = &AddonManager{}
			instance.packageName = packageName
			instance.adapters = make(map[string]IAdapter)
			instance.ipcClient = NewClient(packageName, instance.onMessage)
			instance.running = true
			instance.verbose = true
		},
	)
	return instance
}

func (m *AddonManager) AddAdapters(adapters ...IAdapter) {

	for _, adapter := range adapters {
		m.adapters[adapter.GetID()] = adapter
		adapter.GetID()
		adapter.setManager(m)
		data := make(map[string]any)
		data[Aid] = adapter.GetID()
		data["name"] = adapter.GetName
		data["packageName"] = m.packageName
		m.send(AdapterAddedNotification, data)
	}
}

func (m *AddonManager) handleDeviceAdded(device IDevice) {
	if m.verbose {
		log.Printf("addonManager: handle_device_added: %s", device.GetID())
	}
	data := make(map[string]any)
	data[Aid] = device.GetID()
	description, err := json.Marshal(device)
	if err != nil {
		return
	}
	data["device"] = description
	m.send(DeviceAddedNotification, data)
}

func (m *AddonManager) handleDeviceRemoved(device IDevice) {
	if m.verbose {
		log.Printf("addon manager handle devices added, deviceId:%v\n", device.GetID())
	}
	data := make(map[string]any)
	data[Aid] = device.GetAdapterId()
	data[Did] = device.GetID()

	m.send(AdapterRemoveDeviceResponse, data)
}

func (m *AddonManager) getAdapter(adapterId string) IAdapter {
	adapter, ok := m.adapters[adapterId]
	if !ok {
		return nil
	}
	return adapter
}

func (m *AddonManager) onMessage(data []byte) {

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
		go adapter.StartPairing(timeout)
		return

	case AdapterCancelPairingCommand:
		go adapter.CancelPairing()
		return

		//adapter unload request

	case AdapterUnloadRequest:
		go adapter.Unload()
		unloadFunc := func(proxy *AddonManager, adapter IAdapter) {
			data := make(map[string]any)
			data[Aid] = adapter.GetID()
			proxy.send(AdapterUnloadResponse, data)
		}
		go unloadFunc(m, adapter)
		delete(m.adapters, adapter.GetID())
		break
	}
	var deviceId = json.Get(data, "data", "deviceId").ToString()
	device := adapter.getDevice(deviceId)
	if device == nil {
		return
	}

	switch messageType {
	case AdapterCancelRemoveDeviceCommand:
		adapter := m.adapters[adapterId]
		log.Printf(adapter.GetID())

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
		data[Aid] = adapterId
		data[Did] = device.GetID()
		data["property"] = prop.AsDict()
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
			data[Aid] = adapterId
			data[Did] = deviceId
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

func (m *AddonManager) sendConnectedStateNotification(device *devices.Device, connected bool) {
	data := make(map[string]any)
	data[Aid] = ""
	data[Did] = device.ID
	data["connected"] = connected
	m.send(DeviceConnectedStateNotification, data)
}

func (m *AddonManager) run() {
	m.ipcClient.Register()
}

func (m *AddonManager) send(messageType int, data map[string]any) {
	data[Pid] = m.packageName
	var message = struct {
		MessageType int `json:"messageType"`
		Data        any `json:"data"`
	}{MessageType: messageType, Data: data}
	d, er := json.MarshalIndent(message, "", " ")
	if er != nil {
		log.Fatal(er)
		return
	}
	m.ipcClient.sendMessage(d)
}

func (m *AddonManager) close() {
	m.ipcClient.close()
	m.running = false
}
