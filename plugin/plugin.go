package plugin

import (
	"context"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/ipc"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/util"
	json "github.com/json-iterator/go"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Plugin struct {
	manager     *Manager
	pluginId    string
	exec        string
	execPath    string
	closeChan   chan struct{}
	connection  ipc.Connection
	registered  bool
	logger      logging.Logger
	adapters    sync.Map
	notifiers   sync.Map
	apiHandlers sync.Map
}

func NewPlugin(pluginId string, manager *Manager, log logging.Logger) (plugin *Plugin) {
	plugin = &Plugin{}
	plugin.logger = log
	plugin.registered = false
	plugin.manager = manager
	plugin.pluginId = pluginId
	return
}

func (plugin *Plugin) getId() string {
	return plugin.pluginId
}

func (plugin *Plugin) getName() string {
	return plugin.pluginId
}

func (plugin *Plugin) handleAdapterAdded(adapter *Adapter) {
	plugin.adapters.Store(adapter.getId(), adapter)
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

func (plugin *Plugin) getAdapters() (adapters []*Adapter) {
	plugin.adapters.Range(func(key, value interface{}) bool {
		adapter, ok := value.(*Adapter)
		if ok {
			adapters = append(adapters, adapter)
		}
		return true
	})
	return
}

func (plugin *Plugin) getNotifiers() (notifiers []*Notifier) {
	plugin.notifiers.Range(func(key, value interface{}) bool {
		notifier, ok := value.(*Notifier)
		if ok {
			notifiers = append(notifiers, notifier)
		}
		return true
	})
	return
}

func (plugin *Plugin) getApiHandlers() (apiHandlers []*ApiHandler) {
	plugin.apiHandlers.Range(func(key, value interface{}) bool {
		apiHandler, ok := value.(*ApiHandler)
		if ok {
			apiHandlers = append(apiHandlers, apiHandler)
		}
		return true
	})
	return
}

func (plugin *Plugin) OnMsg(mt messages.MessageType, dt interface{}) {

	data, err := json.Marshal(dt)
	if err != nil {
		plugin.logger.Info("Bad message type")
		return
	}
	switch mt {
	case messages.MessageType_AdapterAddedNotification:
		var message messages.AdapterAddedNotificationJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", util.JsonIndent(dt))
			return
		}
		adapter := NewAdapter(plugin, message.AdapterId, message.Name, message.PackageName, plugin.logger)
		go plugin.handleAdapterAdded(adapter)
		return

	case messages.MessageType_NotifierAddedNotification:
		var message messages.NotifierAddedNotificationJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", util.JsonIndent(dt))
			return
		}
		return
	}

	adapterId := json.Get(data, "adapterId").ToString()
	adapter := plugin.getAdapter(adapterId)
	if adapter == nil {
		plugin.logger.Errorf("plugin message err: %s not found, messageType: %v data: %s", adapterId, mt, util.JsonIndent(dt))
		return
	}
	// adapter handler
	switch mt {
	case messages.MessageType_OutletNotifyResponse:
		var message messages.OutletNotifyRequestJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", util.JsonIndent(dt))
			return
		}
		return

	case messages.MessageType_AdapterUnloadResponse:
		var message messages.AdapterUnloadRequestJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", util.JsonIndent(dt))
			return
		}
		go plugin.handleAdapterUnload(message.AdapterId)
		return

	case messages.MessageType_ApiHandlerApiResponse:
		var message messages.ApiHandlerApiRequestJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", util.JsonIndent(dt))
			return
		}
		return

	case messages.MessageType_ApiHandlerAddedNotification:
		var message messages.ApiHandlerAddedNotificationJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", util.JsonIndent(dt))
			return
		}
		return

	case messages.MessageType_ApiHandlerUnloadResponse:
		var message messages.ApiHandlerUnloadResponseJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", util.JsonIndent(dt))

			return
		}
		return

	case messages.MessageType_PluginUnloadResponse:
		var message messages.PluginUnloadResponseJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", util.JsonIndent(dt))
			return
		}
		plugin.shutdown()
		go plugin.manager.pluginServer.unregisterPlugin(message.PluginId)
		return

	case messages.MessageType_PluginErrorNotification:
		var message messages.PluginErrorNotificationJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", util.JsonIndent(dt))
			return
		}
		return

	case messages.MessageType_DeviceAddedNotification:
		var message messages.DeviceAddedNotificationJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", util.JsonIndent(dt))
			return
		}
		device := newDevice(adapter, message.Device)
		go adapter.handleDeviceAdded(device)
		return
	}

	//device handler
	deviceId := json.Get(data, "deviceId").ToString()
	device := plugin.manager.getDevice(deviceId)
	if device == nil {
		plugin.logger.Errorf("device:%s not found", deviceId)
		return
	}
	switch mt {
	case messages.MessageType_DeviceConnectedStateNotification:
		var message messages.DeviceConnectedStateNotificationJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		device.notifyDeviceConnected(message.Connected)
		return

	case messages.MessageType_DeviceActionStatusNotification:
		var message messages.DeviceActionStatusNotificationJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		device.notifyAction(&addon.ActionDescription{
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
			t, ok := device.requestActionTask.Load(message.ActionId)
			if !ok {
				return
			}
			task := t.(chan bool)
			select {
			case task <- message.Success:
			}
		}()
		return

	case messages.MessageType_DeviceRemoveActionResponse:
		var message messages.DeviceRemoveActionResponseJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
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
			task, ok1 := t.(chan *Device)
			if !ok || task == nil || !ok1 {
				plugin.logger.Errorf("unrecognized message id: %s", message.MessageId)
				return
			}
			if message.Device != nil && message.Success {
				newDev := newDevice(adapter, *message.Device)
				adapter.devices.Store(newDev.GetId(), newDev)
				plugin.manager.devices.Store(newDev.GetId(), newDev)
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
		go device.notifyValueChanged(message.Property)
		return

	case messages.MessageType_AdapterRemoveDeviceResponse:
		var message messages.AdapterRemoveDeviceResponseJsonData
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		go adapter.handleDeviceRemoved(device)
		return

	default:
		plugin.logger.Infof("unknown message mt: %v data: %s", mt, util.JsonIndent(dt))
		return
	}
}

func (plugin *Plugin) send(mt messages.MessageType, data interface{}) {
	err := plugin.connection.WriteMessage(mt, data)
	if err != nil {
		plugin.logger.Infof("plugin %s send err: %s", plugin.pluginId, err.Error())
		return
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

	var syncLog = func(reader io.ReadCloser) {
		for {
			buf := make([]byte, 64)
			for {
				strNum, err := reader.Read(buf)
				if strNum > 0 {
					outputByte := buf[:strNum]
					plugin.logger.Debugf("%s out: %s", plugin.pluginId, string(outputByte))
				}
				if err != nil {
					//读到结尾
					if err == io.EOF || strings.Contains(err.Error(), "file already closed") {
						err = nil
					}
				}
			}
		}
	}

	var cmd *exec.Cmd
	cmd = exec.CommandContext(ctx, commands[0], commands[1:]...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	go syncLog(stdout)
	go syncLog(stderr)
	plugin.logger.Infof("%s start", plugin.pluginId)

	go func() {
		err := cmd.Run()
		if err != nil {
			plugin.logger.Errorf("%s closed: %s", plugin.pluginId, err.Error())
			return
		}
	}()

	for {
		select {
		case <-plugin.closeChan:
			return
		}
	}
}

func (plugin *Plugin) unload() {
	plugin.logger.Info("unloading plugin %s", plugin.pluginId)
	plugin.send(messages.MessageType_PluginUnloadRequest, messages.PluginUnloadRequestJsonData{PluginId: plugin.getId()})
}

func (plugin *Plugin) unloadComponents() {
	adapters := plugin.getAdapters()
	notifiers := plugin.getNotifiers()
	apiHandlers := plugin.getApiHandlers()
	var unloadsFunc []func()
	for _, adapter := range adapters {
		plugin.manager.removeAdapter(adapter)
		for _, device := range adapter.getDevices() {
			plugin.manager.handleDeviceRemoved(device)
		}
		unloadsFunc = append(unloadsFunc, adapter.unload)
	}

	for _, notifier := range notifiers {
		plugin.manager.removeNotifier(notifier.ID)
		unloadsFunc = append(unloadsFunc, notifier.unload)
		for _, device := range notifier.getOutlets() {
			plugin.manager.handleOutletRemoved(device)
		}
	}

	for id, apiHandler := range apiHandlers {
		plugin.manager.removeApiHandler(id)
		unloadsFunc = append(unloadsFunc, apiHandler.unload)
	}
	for _, f := range unloadsFunc {
		f()
	}
	if len(adapters) == 0 && len(notifiers) == 0 && len(apiHandlers) == 0 {
		plugin.unload()
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
