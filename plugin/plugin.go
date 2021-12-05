package plugin

import (
	"context"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/ipc"
	messages "github.com/galenliu/gateway/pkg/ipc_messages"
	"github.com/galenliu/gateway/pkg/logging"
	json "github.com/json-iterator/go"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Plugin struct {
	locker       *sync.Mutex
	pluginId     string
	exec         string
	execPath     string
	restart      bool
	unloading    bool
	closeChan    chan struct{}
	pluginServer *PluginsServer
	clint        ipc.Clint
	logger       logging.Logger
	addonManager *Manager
	adapters     sync.Map
	notifiers    sync.Map
	apiHandlers  sync.Map
	services     sync.Map
}

func NewPlugin(pluginId string, manager *Manager, s *PluginsServer, log logging.Logger) (plugin *Plugin) {
	execPath := s.manager.getAddonPath(pluginId)
	plugin = &Plugin{}
	if execPath != "" {
		plugin.execPath = execPath
	}
	plugin.logger = log
	plugin.locker = new(sync.Mutex)
	plugin.addonManager = manager
	plugin.pluginId = pluginId
	plugin.restart = false
	plugin.pluginServer = s
	return
}

func (plugin *Plugin) getId() string {
	return plugin.pluginId
}

func (plugin *Plugin) getName() string {
	return plugin.pluginId
}

func (plugin *Plugin) getAdapter(adapterId string) *Adapter {
	a, ok := plugin.adapters.Load(adapterId)
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

func (plugin *Plugin) OnMsg(mt messages.MessageType, data interface{}) {

	d, err := json.Marshal(data)
	if err != nil {
		plugin.logger.Errorf("Bad message : %s", err.Error())
		return
	}

	switch mt {
	case messages.MessageType_AdapterAddedNotification:
		var message messages.AdapterAddedNotificationJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		adapter := NewAdapter(plugin, message.AdapterId, message.Name, message.PackageName, plugin.logger)
		plugin.adapters.Store(adapter.id, adapter)
		plugin.addonManager.addAdapter(adapter)
		return

	case messages.MessageType_NotifierAddedNotification:
		var message messages.NotifierAddedNotificationJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
	}

	adapterId := json.Get(d, "adapterId").ToString()
	adapter := plugin.getAdapter(adapterId)
	if adapter == nil {
		plugin.logger.Errorf("adapter not found")
		return
	}
	// adapter handler
	switch mt {
	case messages.MessageType_OutletNotifyResponse:
		var message messages.OutletNotifyRequestJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		break
	case messages.MessageType_AdapterUnloadResponse:
		var message messages.AdapterUnloadRequestJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		break

	case messages.MessageType_ApiHandlerApiResponse:
		var message messages.ApiHandlerApiRequestJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		break

	case messages.MessageType_ApiHandlerAddedNotification:
		var message messages.ApiHandlerAddedNotificationJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		break
	case messages.MessageType_ApiHandlerUnloadResponse:
		var message messages.ApiHandlerUnloadResponseJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		break

	case messages.MessageType_PluginUnloadRequest:
		var message messages.PluginUnloadRequestJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		plugin.pluginServer.unregisterPlugin(message.PluginId)
		break

	case messages.MessageType_PluginErrorNotification:
		var message messages.PluginErrorNotificationJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		break

	case messages.MessageType_DeviceAddedNotification:
		var message messages.DeviceAddedNotificationJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		device := newDevice(adapter, message.Device)
		adapter.handleDeviceAdded(device)
		return
	}

	//device handler
	deviceId := json.Get(d, "deviceId").ToString()
	device := plugin.pluginServer.manager.getDevice(deviceId)
	if device == nil {
		plugin.logger.Errorf("device:%s not found", deviceId)
		return
	}
	switch mt {
	case messages.MessageType_DeviceConnectedStateNotification:
		var message messages.DeviceConnectedStateNotificationJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		device.notifyDeviceConnected(message.Connected)
		break

	case messages.MessageType_DeviceActionStatusNotification:
		var message messages.DeviceActionStatusNotificationJsonData
		err := json.Unmarshal(d, &message)
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
		break

	case messages.MessageType_DeviceRequestActionResponse:
		var message messages.DeviceRequestActionResponseJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		task, ok := device.requestActionTask[message.ActionId]
		if !ok {
			return
		}
		select {
		case task <- message.Success:
		}
		break
	case messages.MessageType_DeviceRemoveActionResponse:
		var message messages.DeviceRemoveActionResponseJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		break
	case messages.MessageType_DeviceSetCredentialsResponse:
		var message messages.DeviceSetCredentialsRequestJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		break
	case messages.MessageType_DevicePropertyChangedNotification:
		var message messages.DevicePropertyChangedNotificationJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		if message.Property.Name == nil || message.Property.Value == nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		var readOnly = false
		if message.Property.ReadOnly == nil {
			message.Property.ReadOnly = &readOnly
		}
		device.notifyValueChanged(message.Property)

	case messages.MessageType_AdapterRemoveDeviceResponse:
		var message messages.AdapterRemoveDeviceResponseJsonData
		err := json.Unmarshal(d, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : %s", err.Error())
			return
		}
		adapter.handleDeviceRemoved(device)
		return
	default:
		plugin.logger.Infof("unknown message : %v", mt)
	}
	return
}

func (plugin *Plugin) sendMsg(mt messages.MessageType, message map[string]interface{}) {
	message["pluginId"] = plugin.pluginId
	err := plugin.clint.WriteMessage(mt, message)
	if err != nil {
		plugin.logger.Infof("plugin %s send err: %s", plugin.pluginId, err.Error())
		return
	}
}

func (plugin *Plugin) send(mt messages.MessageType, data interface{}) {
	err := plugin.clint.WriteMessage(mt, data)
	if err != nil {
		plugin.logger.Infof("plugin %s send err: %s", plugin.pluginId, err.Error())
		return
	}
}

func (plugin *Plugin) run() {
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
		plugin.logger.Warningf("plugin %s not run,only support plugin with python lang now", plugin.pluginId)
		return
	}
	commands := strings.Split(command, " ")
	var syncLog = func(reader io.ReadCloser) {
		ctx1, f := context.WithCancel(ctx)
		for {
			select {
			case <-ctx1.Done():
				f()
				return
			default:
				buf := make([]byte, 1024)
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

	}
	var cmd *exec.Cmd
	var args = commands[1:]
	if len(args) > 0 {
		cmd = exec.CommandContext(ctx, commands[0], args...)
	} else {
		cmd = exec.CommandContext(ctx, commands[0])
	}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	go syncLog(stdout)
	go syncLog(stderr)
	plugin.logger.Infof("%s run", plugin.pluginId)

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
	plugin.restart = false
	plugin.unloading = true
	plugin.sendMsg(messages.MessageType_PluginUnloadRequest, map[string]interface{}{})
}

func (plugin *Plugin) unloadComponents() {
	adapters := plugin.getAdapters()
	notifiers := plugin.getNotifiers()
	apiHandlers := plugin.getApiHandlers()
	var unloadsFunc []func()
	for _, adapter := range adapters {
		plugin.addonManager.removeAdapter(adapter)
		for _, device := range adapter.getDevices() {
			plugin.addonManager.handleDeviceRemoved(device)
		}
		unloadsFunc = append(unloadsFunc, adapter.unload)
	}

	for _, notifier := range notifiers {
		plugin.addonManager.removeNotifier(notifier.ID)
		unloadsFunc = append(unloadsFunc, notifier.unload)
		for _, device := range notifier.getOutlets() {
			plugin.addonManager.handleOutletRemoved(device)
		}
	}

	for id, apiHandler := range apiHandlers {
		plugin.addonManager.removeApiHandler(id)
		unloadsFunc = append(unloadsFunc, apiHandler.unload)
	}
	for _, f := range unloadsFunc {
		f()
	}
	if len(adapters) == 0 && len(notifiers) == 0 && len(apiHandlers) == 0 {
		plugin.unload()
	}
}

// kill
//  @Description:  kill the process
func (plugin *Plugin) kill() {
	select {
	case plugin.closeChan <- struct{}{}:
	}
}
