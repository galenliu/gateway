package proxy

import (
	"context"
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/galenliu/gateway/pkg/addon/manager"
	"github.com/galenliu/gateway/pkg/addon/properties"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	json "github.com/json-iterator/go"
	"log"
	"net/netip"
	"time"
)

type Manager struct {
	*manager.Manager
	pluginId    string
	verbose     bool
	running     bool
	registered  bool
	Done        chan struct{}
	ctx         context.Context
	sendChan    chan []byte
	userProfile *messages.PluginRegisterResponseJsonDataUserProfile
	preferences *messages.PluginRegisterResponseJsonDataPreferences
}

func NewAddonManager(ctx context.Context, pluginId string) (*Manager, error) {

	instance := &Manager{}
	instance.Manager = manager.NewManager()
	instance.pluginId = pluginId
	instance.verbose = true
	instance.Done = make(chan struct{})
	//instance.ipcClient = NewClient(instance.ctx, instance, "9500")
	//if instance.ipcClient != nil {
	//	instance.running = true
	//	instance.register()
	//}
	addr, err := netip.ParseAddr("127.0.0.1")
	if err != nil {
		return nil, err
	}
	host := netip.AddrPortFrom(addr, 9500)
	instance.sendChan = make(chan []byte, 256)
	client := &Client{
		send: instance.sendChan,
		host: &host,
	}
	err = instance.register(client)
	if err != nil {
		return nil, err
	}
	readChan := make(chan []byte, 256)
	go client.readPump(ctx, readChan)
	go client.writePump()
	go func(ctx context.Context) {
		for {
			select {
			case data := <-readChan:
				instance.OnMessage(data)
			case <-ctx.Done():
				log.Printf("client exit")
			}
		}
	}(ctx)
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
			time.AfterFunc(500*time.Millisecond, func() { m.close() })
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
	bytes, err := json.Marshal(message)
	if err != nil {
		log.Printf(err.Error())
	}
	select {
	case m.sendChan <- bytes:
	default:
		log.Printf("send channel is full")
	}
}

func (m *Manager) register(client *Client) error {

	conn, err := client.Dial()
	if err != nil {
		return err
	}
	//err = conn.SetWriteDeadline(time.Now().Add(writeWait))
	//if err != nil {
	//	return err
	//}
	//w, err := client.conn.NextWriter(websocket.TextMessage)
	//if err != nil {
	//	return err
	//}
	data, err := json.Marshal(messages.PluginRegisterRequestJson{
		Data:        messages.PluginRegisterRequestJsonData{PluginId: m.pluginId},
		MessageType: int(messages.MessageType_PluginRegisterRequest),
	})
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	//_, err = w.Write(data)
	//if err != nil {
	//	return err
	//}
	//err = client.conn.SetReadDeadline(time.Now().Add(pongWait))
	//if err != nil {
	//	return err
	//}
	var message messages.PluginRegisterResponseJson
	_, data, err = conn.ReadMessage()
	err = json.Unmarshal(data, &message)
	if err != nil {
		return err
	}
	m.registered = true
	m.preferences = &message.Data.Preferences
	m.userProfile = &message.Data.UserProfile
	return nil
}

func (m *Manager) close() {
	select {
	case m.Done <- struct{}{}:
	}
	log.Printf("plugin %s unloaded", m.pluginId)
}

func (m *Manager) getPluginId() string {
	return m.pluginId
}

func (m *Manager) getUserProfile() *messages.PluginRegisterResponseJsonDataUserProfile {
	return m.userProfile
}
