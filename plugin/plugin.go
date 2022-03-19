package plugin

import (
	"context"
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/galenliu/gateway/pkg/addon/actions"
	"github.com/galenliu/gateway/pkg/addon/properties"
	"github.com/galenliu/gateway/pkg/bus/topic"
	"github.com/galenliu/gateway/pkg/ipc"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	json "github.com/json-iterator/go"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Plugin struct {
	manager      *Manager
	pluginId     string
	exec         string
	execPath     string
	closeChan    chan struct{}
	connection   ipc.Connection
	registered   bool
	logger       logging.Logger
	sendChan     chan []byte
	eventHandler map[string]func()
	adapters     sync.Map
	notifiers    sync.Map
	apiHandlers  sync.Map
}

func NewPlugin(pluginId string, manager *Manager, log logging.Logger) (plugin *Plugin) {
	plugin = &Plugin{}
	plugin.logger = log
	plugin.registered = false
	plugin.manager = manager
	plugin.pluginId = pluginId
	plugin.sendChan = make(chan []byte, 24)
	plugin.eventHandler = make(map[string]func(), 0)
	return
}

func (plugin *Plugin) getId() string {
	return plugin.pluginId
}

func (plugin *Plugin) getName() string {
	return plugin.pluginId
}

func (plugin *Plugin) handleAdapterAdded(adapter *Adapter) {
	plugin.adapters.Store(adapter.GetId(), adapter)
	plugin.manager.handleAdapterAdded(adapter)
}

func (plugin *Plugin) handleAdapterUnload(adapterId string) {
	for _, adapter := range plugin.getAdapters() {
		adapter.unload()
	}
	plugin.adapters.Delete(adapterId)
	plugin.manager.handleAdapterUnload(adapterId)
}

func (plugin *Plugin) getAdapter(adapterId string) *Adapter {
	a, ok := plugin.adapters.Load(adapterId)
	if !ok {
		return nil
	}
	adapter, ok := a.(*Adapter)
	if !ok {
		return nil
	}
	return adapter
}

func (plugin *Plugin) handleWs(ws *websocket.Conn) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go func(ctx context.Context) {
		for {
			select {
			case d := <-plugin.sendChan:
				err := ws.WriteMessage(websocket.TextMessage, d)
				if err != nil {
					plugin.logger.Infof("websocket write error: %v", err)
				}
			case <-ctx.Done():
				plugin.logger.Infof("websocket stop write")
				return
			}
		}
	}(ctx)

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if websocket.IsCloseError(err) {
			return
		}
		plugin.OnMsg(data)
	}

}

func (plugin *Plugin) getAdapters() (adapters []*Adapter) {
	adapters = make([]*Adapter, 0)
	plugin.adapters.Range(func(key, value any) bool {
		adapter, ok := value.(*Adapter)
		if ok {
			adapters = append(adapters, adapter)
		}
		return true
	})
	return
}

func (plugin *Plugin) getNotifiers() (notifiers []*Notifier) {
	notifiers = make([]*Notifier, 0)
	plugin.notifiers.Range(func(key, value any) bool {
		notifier, ok := value.(*Notifier)
		if ok {
			notifiers = append(notifiers, notifier)
		}
		return true
	})
	return
}

func (plugin *Plugin) getApiHandlers() (apiHandlers []*ApiHandler) {
	apiHandlers = make([]*ApiHandler, 0)
	plugin.apiHandlers.Range(func(key, value any) bool {
		apiHandler, ok := value.(*ApiHandler)
		if ok {
			apiHandlers = append(apiHandlers, apiHandler)
		}
		return true
	})
	return
}

func (plugin *Plugin) OnMsg(bt []byte) {
	type msg struct {
		MessageType messages.MessageType `json:"messageType"`
		Data        any                  `json:"data"`
	}
	var m msg
	err := json.Unmarshal(bt, &m)
	if err != nil {
		plugin.logger.Info("Bad message")
		return
	}
	data, err := json.Marshal(m.Data)
	if err != nil {
		plugin.logger.Info("Bad message")
		return
	}

	//首先处理Adapter注册
	{
		switch m.MessageType {
		case messages.MessageType_AdapterAddedNotification:
			var message messages.AdapterAddedNotificationJsonData
			err := json.Unmarshal(data, &message)
			if err != nil {
				plugin.logger.Errorf("Bad message %s", messages.MessageType_AdapterAddedNotification)
				return
			}
			adapter := NewAdapter(message.AdapterId, message.Name, message.PackageName, plugin)

			send := func(msg topic.ThingAddedMessage) {
				var device messages.DeviceWithoutId
				err := json.Unmarshal(msg.Data, &device)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				adapter.Send(messages.MessageType_DeviceSavedNotification, messages.DeviceSavedNotificationJsonData{
					AdapterId: adapter.GetId(),
					Device:    device,
					DeviceId:  msg.ThingId,
					PluginId:  plugin.getId(),
				})
			}

			if plugin.manager.thingContainer != nil {
				for _, t := range plugin.manager.thingContainer.GetThings() {
					send(topic.ThingAddedMessage{
						ThingId: t.GetId(),
						Data:    []byte(util.JsonIndent(t)),
					})
				}
			}
			adapter.eventHandler[string(topic.ThingAdded)] = plugin.manager.thingContainer.Subscribe(topic.ThingAdded, send)

			go plugin.handleAdapterAdded(adapter)
			return

		case messages.MessageType_NotifierAddedNotification:
			var message messages.NotifierAddedNotificationJsonData
			err := json.Unmarshal(data, &message)
			if err != nil {
				plugin.logger.Errorf("Bad message : %s", messages.MessageType_NotifierAddedNotification)
				return
			}
			return
		}
	}

	//处理adapter消息，如果adapter未注册，则丢弃消息
	adapterId := json.Get(data, "adapterId").ToString()
	adapter := plugin.getAdapter(adapterId)
	if adapter == nil {
		plugin.logger.Errorf("can not find adapter %s", adapterId)
		return
	}
	switch m.MessageType {
	case messages.MessageType_OutletNotifyResponse:
		var message messages.OutletNotifyRequestJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", messages.MessageType_OutletNotifyResponse)
			return
		}
		return

	case messages.MessageType_AdapterUnloadResponse:
		var message messages.AdapterUnloadRequestJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", messages.MessageType_AdapterUnloadResponse)
			return
		}
		go plugin.handleAdapterUnload(message.AdapterId)
		return

	case messages.MessageType_ApiHandlerApiResponse:
		var message messages.ApiHandlerApiRequestJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", messages.MessageType_ApiHandlerApiResponse)
			return
		}
		return

	case messages.MessageType_ApiHandlerAddedNotification:
		var message messages.ApiHandlerAddedNotificationJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", messages.MessageType_ApiHandlerAddedNotification)
			return
		}
		return

	case messages.MessageType_ApiHandlerUnloadResponse:
		var message messages.ApiHandlerUnloadResponseJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", m.MessageType)
			return
		}
		return

	case messages.MessageType_PluginUnloadResponse:
		var message messages.PluginUnloadResponseJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", m.MessageType)
			return
		}
		plugin.shutdown()
		go plugin.manager.pluginServer.unregisterPlugin(message.PluginId)
		return

	case messages.MessageType_PluginErrorNotification:
		var message messages.PluginErrorNotificationJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", m.MessageType)
			return
		}
		return

	case messages.MessageType_DeviceAddedNotification:
		var message messages.DeviceAddedNotificationJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", m.MessageType)
			return
		}
		device := newDeviceFromMessage(message.Device)
		go adapter.handleDeviceAdded(device)
		return
	}

	//处理device消息，如果device没有注册，则丢弃消息
	deviceId := json.Get(data, "deviceId").ToString()
	dev := plugin.manager.getDevice(deviceId)
	if dev == nil {
		plugin.logger.Errorf("device:%s not found", deviceId)
		return
	}
	switch m.MessageType {
	case messages.MessageType_DeviceConnectedStateNotification:
		var message messages.DeviceConnectedStateNotificationJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		dev.notifyDeviceConnected(message.Connected)
		return

	case messages.MessageType_DeviceActionStatusNotification:
		var message messages.DeviceActionStatusNotificationJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		dev.actionNotify(actions.ActionDescription{
			Id:            message.Action.Id,
			Name:          message.Action.Name,
			Input:         message.Action.Input,
			Status:        message.Action.Status,
			TimeRequested: message.Action.TimeRequested,
			TimeCompleted: func() string {
				if message.Action.TimeCompleted == nil {
					return ""
				}
				return *message.Action.TimeCompleted
			}(),
		})
		return

	case messages.MessageType_DeviceRequestActionResponse:
		var message messages.DeviceRequestActionResponseJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		go func() {
			t, ok := dev.requestActionTask.Load(message.ActionId)
			if !ok {
				return
			}
			task, ok := t.(chan bool)
			if ok {
				select {
				case task <- message.Success:
				}
			}
			dev.requestActionTask.Delete(message.ActionId)
		}()
		return

	case messages.MessageType_DeviceRemoveActionResponse:
		var message messages.DeviceRemoveActionResponseJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		go func() {
			t, ok := dev.removeActionTask.Load(message.ActionId)
			if !ok {
				return
			}
			task, ok := t.(chan bool)
			if ok {
				select {
				case task <- message.Success:
				}
			}
			dev.removeActionTask.Delete(message.ActionId)
		}()
		return

	case messages.MessageType_DeviceSetCredentialsResponse:
		var message messages.DeviceSetCredentialsResponseJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		go func() {
			t, ok := adapter.setCredentialsTask.Load(message.MessageId)
			task, ok1 := t.(chan *device)
			if !ok || task == nil || !ok1 {
				plugin.logger.Errorf("unrecognized message id: %s", message.MessageId)
				return
			}
			if message.Device != nil && message.Success {
				newDev := newDeviceFromMessage(*message.Device)
				adapter.handleDeviceAdded(newDev)
				select {
				case task <- newDev:
				}
			}
		}()
		return

	case messages.MessageType_DevicePropertyChangedNotification:
		var message messages.DevicePropertyChangedNotificationJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		if message.Property.Name == nil || message.Property.Value == nil {
			plugin.logger.Errorf("Bad message")
			return
		}
		go dev.onPropertyChanged(properties.PropertyDescription{
			Name:        util.GetValueFromPointer(message.Property.Name),
			AtType:      *message.Property.AtType,
			Title:       util.GetValueFromPointer(message.Property.Title),
			Type:        message.Property.Type,
			Unit:        util.GetValueFromPointer(message.Property.Unit),
			Description: util.GetValueFromPointer(message.Property.Description),
			Minimum:     util.GetValueFromPointer(message.Property.Minimum),
			Maximum:     util.GetValueFromPointer(message.Property.Maximum),
			Enum:        message.Property.Enum,
			ReadOnly:    util.GetValueFromPointer(message.Property.ReadOnly),
			MultipleOf:  util.GetValueFromPointer(message.Property.MultipleOf),
			Links:       nil,
			Value:       message.Property.Value,
		})
		return

	case messages.MessageType_AdapterRemoveDeviceResponse:
		var message messages.AdapterRemoveDeviceResponseJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		go adapter.handleDeviceRemoved(dev)
		return

	default:
		plugin.logger.Infof("unknown message:", m.MessageType)
		return
	}
}

func (plugin *Plugin) send(mt messages.MessageType, d any) {
	type msg struct {
		MessageType messages.MessageType
		Data        any
	}
	m := msg{
		MessageType: mt,
		Data:        d,
	}
	data, err := json.Marshal(m)
	if err != nil {
		plugin.logger.Infof("bad message")
		return
	}
	if plugin.registered == true {
		select {
		case plugin.sendChan <- data:
		}
	}
}

func (plugin *Plugin) start() {

	plugin.closeChan = make(chan struct{})
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	plugin.exec = strings.Replace(plugin.exec, "\\", string(os.PathSeparator), -1)
	plugin.exec = strings.Replace(plugin.exec, "/", string(os.PathSeparator), -1)
	plugin.execPath = strings.Replace(plugin.execPath, "\\", string(os.PathSeparator), -1)
	plugin.execPath = strings.Replace(plugin.execPath, "/", string(os.PathSeparator), -1)
	command := strings.Replace(plugin.exec, "{path}", plugin.execPath, 1)
	//command = strings.Replace(command, "{nodeLoader}", configs.GetNodeLoader(), 1)
	if !strings.HasPrefix(command, "python") {
		plugin.logger.Warningf("plugin %s only support with python", plugin.getId())
		return
	}
	commands := strings.Split(command, " ")
	//
	//var syncLog = func(reader io.ReadCloser) {
	//	for {
	//		buf := make([]byte, 64)
	//		for {
	//			strNum, err := reader.Read(buf)
	//			if strNum > 0 {
	//				outputByte := buf[:strNum]
	//				plugin.logger.Debugf("%s out: %s", plugin.pluginId, string(outputByte))
	//			}
	//			if err != nil {
	//				//读到结尾
	//				if err == io.EOF || strings.Contains(err.Error(), "file already closed") {
	//					err = nil
	//				}
	//			}
	//		}
	//	}
	//}
	//
	var cmd *exec.Cmd
	cmd = exec.CommandContext(ctx, commands[0], commands[1:]...)
	//stdout, _ := cmd.StdoutPipe()
	//stderr, _ := cmd.StderrPipe()
	//
	//go syncLog(stdout)
	//go syncLog(stderr)
	//plugin.logger.Infof("%s start", plugin.pluginId)
	//
	go func() {
		err := cmd.Run()
		if err != nil {
			plugin.logger.Errorf("%s closed: %s", plugin.pluginId, err.Error())
			return
		}
	}()
	//
	//for {
	//	select {
	//	case <-plugin.closeChan:
	//		return
	//	}
	//}
}

func (plugin *Plugin) notifyUnload() {
	plugin.logger.Infof("unloading plugin %s", plugin.getId())
	plugin.send(messages.MessageType_PluginUnloadRequest, messages.PluginUnloadRequestJsonData{PluginId: plugin.getId()})
}

func (plugin *Plugin) unloadComponents() {
	adapters := plugin.getAdapters()
	notifiers := plugin.getNotifiers()
	apiHandlers := plugin.getApiHandlers()
	var unloadsFunc = make([]func(), 0)
	if adapters != nil && len(adapters) > 0 {
		for _, adapter := range adapters {
			plugin.manager.removeAdapter(adapter)
			for _, device := range adapter.getDevices() {
				plugin.manager.handleDeviceRemoved(device)
			}
			unloadsFunc = append(unloadsFunc, adapter.unload)
		}
	}

	if notifiers != nil && len(notifiers) > 0 {
		for _, notifier := range notifiers {
			plugin.manager.removeNotifier(notifier.ID)
			unloadsFunc = append(unloadsFunc, notifier.unload)
			for _, device := range notifier.getOutlets() {
				plugin.manager.handleOutletRemoved(device)
			}
		}
	}

	if apiHandlers != nil && len(apiHandlers) > 0 {
		for id, apiHandler := range apiHandlers {
			plugin.manager.removeApiHandler(id)
			unloadsFunc = append(unloadsFunc, apiHandler.unload)
		}
	}

	if unloadsFunc != nil && len(unloadsFunc) > 0 {
		for _, f := range unloadsFunc {
			go f()
		}
	}
	if len(adapters) == 0 && len(notifiers) == 0 && len(apiHandlers) == 0 {
		go plugin.notifyUnload()
	}
}

func (plugin *Plugin) register(connection ipc.Connection) {
	plugin.registered = true
	plugin.connection = connection
	connection.Register(plugin.getId())
}

// closed the plugin connection
func (plugin *Plugin) shutdown() {
	select {
	case plugin.closeChan <- struct{}{}:
	}
}
