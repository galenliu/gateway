package proxy

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/addon/manager"
	"github.com/galenliu/gateway/pkg/addon/properties"
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
	ctx         context.Context
	userProfile *messages.PluginRegisterResponseJsonDataUserProfile
	preferences *messages.PluginRegisterResponseJsonDataPreferences
}

var once sync.Once
var instance *Manager

func NewAddonManager(pluginId string) (*Manager, error) {
	once.Do(
		func() {
			instance = &Manager{}
			instance.ctx = context.Background()
			instance.Manager = manager.NewManager()
			instance.pluginId = pluginId
			instance.verbose = true
			instance.registered = false
			instance.ipcClient = NewClient(instance.ctx, instance, "9500")
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

func (m *Manager) AddAdapters(adapters ...AdapterProxy) {
	for _, adapter := range adapters {
		adapter.registered(m)
		m.Manager.AddAdapter(adapter)
		m.send(messages.MessageType_AdapterAddedNotification, messages.AdapterAddedNotificationJsonData{
			AdapterId:   adapter.GetId(),
			Name:        adapter.GetName(),
			PackageName: adapter.GetPackageName(),
			PluginId:    m.pluginId,
		})
	}
}

func (m *Manager) handleDeviceAdded(device DeviceProxy) {
	if m.verbose {
		fmt.Printf("manager device_added: %s \t\n", device.GetId())
	}
	m.send(messages.MessageType_DeviceAddedNotification, messages.DeviceAddedNotificationJsonData{
		AdapterId: device.GetAdapter().GetId(),
		Device:    *device.ToMessage(),
		PluginId:  m.pluginId,
	})
}

func (m *Manager) handleDeviceRemoved(device DeviceProxy) {
	if m.verbose {
		fmt.Printf("addon manager handle devices added, deviceId:%v\n", device.GetId())
	}
	m.send(messages.MessageType_AdapterRemoveDeviceResponse, messages.AdapterRemoveDeviceResponseJsonData{
		AdapterId: device.GetAdapter().GetId(),
		DeviceId:  device.GetId(),
		PluginId:  m.pluginId,
	})
}

func (m *Manager) HandleAdapterRemoved(id string) {
	m.Manager.RemoveAdapter(id)
}

func (m *Manager) GetAdapter(adapterId string) AdapterProxy {
	adapter := m.Manager.GetAdapter(adapterId)
	if adapter != nil {
		adp, ok := adapter.(AdapterProxy)
		if ok {
			return adp
		}
	}
	return nil
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

func (m *Manager) getDevice(deviceId string) DeviceProxy {
	adapter := m.Manager.GetDevice(deviceId)
	if adapter != nil {
		dev, ok := adapter.(DeviceProxy)
		if ok {
			return dev
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
		m.send(messages.MessageType_PluginUnloadResponse, messages.PluginUnloadResponseJsonData{PluginId: msg.PluginId})
		m.running = false
		go func() {
			time.AfterFunc(500*time.Millisecond, func() { m.Close() })
		}()
		return
	}

	adapterId := dataNode.Get("adapterId").ToString()
	adapter := m.getAdapter(adapterId)
	if adapter == nil {
		fmt.Printf("can not found adapter(%s) \t\n", adapterId)
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
		go adapter.unload()
		unloadFunc := func() {
			m.send(messages.MessageType_AdapterUnloadResponse, messages.AdapterUnloadResponseJsonData{
				AdapterId: adapter.GetId(),
				PluginId:  m.pluginId,
			})
			m.HandleAdapterRemoved(adapter.GetId())
		}
		go unloadFunc()
		break

	case messages.MessageType_DeviceSavedNotification:
		var msg messages.DeviceSavedNotificationJsonData
		dataNode.ToVal(&msg)
		if dataNode.LastError() != nil {
			fmt.Printf("message unmarshal err:%s", dataNode.LastError().Error())
			return
		}
		go adapter.HandleDeviceSaved(msg)
		return
	}

	deviceId := dataNode.Get("deviceId").ToString()
	device := adapter.getDevice(deviceId)
	if device == nil {
		fmt.Printf("manager Onmessage: device %s not found \t\n", deviceId)
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

	case messages.MessageType_AdapterRemoveDeviceRequest:
		var msg messages.AdapterRemoveDeviceRequestJsonData
		dataNode.ToVal(&msg)
		if dataNode.LastError() != nil {
			fmt.Printf("message unmarshal err:%s", dataNode.LastError().Error())
			return
		}
		go func() {
			adapter.HandleDeviceRemoved(device)
			adapter.RemoveDevice(deviceId)
		}()
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
		go func() {
			{
				{
					switch prop.GetType() {
					case properties.TypeBoolean:
						p, ok := prop.(properties.BooleanEntity)
						if ok {
							b := p.CheckValue(msg.PropertyValue)
							if b {
								err := p.TurnOn()
								if err != nil {
									fmt.Printf(err.Error())
								}
							} else {
								err := p.TurnOff()
								if err != nil {
									fmt.Printf(err.Error())
								}
							}
						} else {
							p, ok := prop.(properties.Entity)
							if ok {
								err := p.SetPropertyValue(msg.PropertyValue)
								if err != nil {
									fmt.Printf(err.Error())
								}
							}
						}
					case properties.TypeInteger:
						p, ok := prop.(properties.IntegerEntity)
						if ok {
							value := p.CheckValue(msg.PropertyValue)
							err := p.SetValue(value)
							if err != nil {
								fmt.Printf(err.Error())
							}
						} else {
							p, ok := prop.(properties.Entity)
							if ok {
								err := p.SetPropertyValue(msg.PropertyValue)
								if err != nil {
									fmt.Printf(err.Error())
								}
							}
						}
					case properties.TypeNumber:
						p, ok := prop.(properties.NumberEntity)
						if ok {
							value := p.CheckValue(msg.PropertyValue)
							err := p.SetValue(value)
							if err != nil {
								fmt.Printf(err.Error())
							}
						} else {
							p, ok := prop.(properties.Entity)
							if ok {
								err := p.SetPropertyValue(msg.PropertyValue)
								if err != nil {
									fmt.Printf(err.Error())
								}
							}
						}
					case properties.TypeString:
						p, ok := prop.(properties.StringEntity)
						if ok {
							value := p.CheckValue(msg.PropertyValue)
							err := p.SetValue(value)
							if err != nil {
								fmt.Printf(err.Error())
							}
						} else {
							p, ok := prop.(properties.Entity)
							if ok {
								err := p.SetPropertyValue(msg.PropertyValue)
								if err != nil {
									fmt.Printf(err.Error())
								}
							}
						}
					default:
						p, ok := prop.(properties.Entity)
						if ok {
							err := p.SetPropertyValue(msg.PropertyValue)
							if err != nil {
								fmt.Printf(err.Error())
							}
						}
					}
				}

			}
		}()
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
			m.send(messages.MessageType_DeviceSetPinResponse, messages.DeviceSetPinResponseJsonData{
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
			m.send(messages.MessageType_DeviceSetCredentialsResponse, messages.DeviceSetCredentialsResponseJsonData{
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
	m.send(messages.MessageType_DeviceConnectedStateNotification, messages.DeviceConnectedStateNotificationJsonData{
		AdapterId: device.GetAdapter().GetId(),
		Connected: connected,
		DeviceId:  device.GetId(),
		PluginId:  m.getPluginId(),
	})
}

func (m *Manager) send(messageType messages.MessageType, data any) {
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
	m.send(messages.MessageType_PluginRegisterRequest, messages.PluginRegisterRequestJsonData{PluginId: m.pluginId})
	time.Sleep(time.Duration(1000) * time.Millisecond)
}

func (m *Manager) Close() {
	m.ctx.Done()
	m.running = false
}

func (m *Manager) IsRunning() bool {
	if len(m.GetAdapters()) > 0 && m.running {
		return true
	}
	return false
}

func (m *Manager) getPluginId() string {
	return m.pluginId
}

func (m *Manager) getUserProfile() *messages.PluginRegisterResponseJsonDataUserProfile {
	return m.userProfile
}
