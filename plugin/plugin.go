package plugin

import (
	"context"
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/addon"
	"github.com/galenliu/gateway/pkg/ipc"
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

func (plugin *Plugin) getService(id string) *Service {
	a, ok := plugin.services.Load(id)
	s, ok := a.(*Service)
	if !ok {
		return nil
	}
	return s
}

func (plugin *Plugin) getServices() (services []*Service) {
	plugin.services.Range(func(key, value interface{}) bool {
		adapter, ok := value.(*Service)
		if ok {
			services = append(services, adapter)
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

func (plugin *Plugin) OnMsg(messageType rpc.MessageType, data []byte) {

	switch messageType {
	case rpc.MessageType_AdapterAddedNotification:
		var message rpc.AdapterAddedNotificationMessage_Data
		err := json.Unmarshal(data, &message)
		if err != nil {
			plugin.logger.Errorf("Bad message : MessageType_AdapterAddedNotification")
			return
		}
		adapter := NewAdapter(plugin, message.AdapterId, message.Name, message.PackageName, plugin.logger)
		plugin.adapters.Store(adapter.ID, adapter)
		plugin.addonManager.addAdapter(adapter)
		return

	case rpc.MessageType_NotifierAddedNotification:
		return

	}
	adapterId := json.Get(data, "adapterId").ToString()
	adapter := plugin.getAdapter(adapterId)
	if adapter == nil {
		plugin.logger.Errorf("adapter not found")
		return
	}

	// adapter handler
	switch messageType {

	case rpc.MessageType_OutletNotifyResponse:
		break
	case rpc.MessageType_AdapterUnloadResponse:
		break

	case rpc.MessageType_ApiHandlerApiResponse:
		break

	case rpc.MessageType_ApiHandlerAddedNotification:
		return
	case rpc.MessageType_ApiHandlerUnloadResponse:
		return
	case rpc.MessageType_PluginUnloadRequest:
		plugin.pluginServer.unregisterPlugin(plugin.pluginId)
		return
	case rpc.MessageType_PluginErrorNotification:
		return
	case rpc.MessageType_DeviceAddedNotification:
		var deviceDescription addon.Device
		json.Get(data, "device").ToVal(&deviceDescription)
		if &deviceDescription == nil {
			plugin.logger.Errorf("deivce added notification bad message")
			return
		}
		var device = &Device{
			adapter: adapter,
			Device:  &deviceDescription,
		}
		adapter.handleDeviceAdded(device)
		return
	}

	// services message
	{
		switch messageType {
		case rpc.MessageType_ServiceAddedNotification:
			var d rpc.ServiceAddedNotificationMessage_Data
			err := json.Unmarshal(data, &d)
			if err != nil {
				return
			}
			newService := NewService(plugin, plugin.addonManager.Eventbus, d.ServiceId, d.Name)
			plugin.services.Store(newService.ID, newService)
			plugin.addonManager.addService(newService)

		case rpc.MessageType_ServiceGetThingsRequest:
			things := plugin.pluginServer.manager.container.GetThings()
			bt, err := json.Marshal(things)
			if err != nil {
				return
			}
			var d = make(map[string]interface{})
			d["things"] = bt
			adapter.sendMsg(rpc.MessageType_ServiceGetThingsResponse, d)

		case rpc.MessageType_ServiceGetThingRequest:
			id := json.Get(data, "thingId").ToString()
			thing := plugin.pluginServer.manager.container.GetThing(id)
			bt, err := json.Marshal(thing)
			if err != nil {
				return
			}
			var d = make(map[string]interface{})
			d["thing"] = bt
			adapter.sendMsg(rpc.MessageType_ServiceGetThingResponse, d)

		case rpc.MessageType_ServiceSetPropertyValueRequest:
			var message rpc.ServiceSetPropertyValueRequestMessage_Data
			err := json.Unmarshal(data, &message)
			if err != nil {
				return
			}
			_, err = adapter.plugin.pluginServer.manager.SetPropertyValue(message.ThingId, message.PropertyName, message.Value)
			if err != nil {
				return
			}
		}
	}

	// addon handler
	{
		switch messageType {
		case rpc.MessageType_AdapterUnloadResponse:
			return
		case rpc.MessageType_NotifierUnloadResponse:
			return

		case rpc.MessageType_OutletAddedNotification:
			return
		case rpc.MessageType_OutletRemovedNotification:
			return

		case rpc.MessageType_AdapterPairingPromptNotification:
			return

		case rpc.MessageType_AdapterUnpairingPromptNotification:
			return
		case rpc.MessageType_MockAdapterClearStateResponse:
			return

		case rpc.MessageType_MockAdapterRemoveDeviceResponse:
			return

		}
	}

	//device handler
	device := plugin.pluginServer.manager.getDevice(json.Get(data, "deviceId").ToString())
	if device == nil {
		return
	}
	switch messageType {
	case rpc.MessageType_DeviceRequestActionResponse:
		break
	case rpc.MessageType_DeviceRemoveActionResponse:
		break
	case rpc.MessageType_DeviceSetCredentialsResponse:
		break
	case rpc.MessageType_DevicePropertyChangedNotification:
		var property addon.Property
		json.Get(data, "property").ToVal(&property)
		if &property == nil {
			plugin.logger.Errorf("device property changed notification bad message")
			return
		}
		device.notifyValueChanged(&property)

	case rpc.MessageType_AdapterRemoveDeviceResponse:
		adapter.handleDeviceRemoved(device)
		return
	default:
		plugin.logger.Infof("unkown message : %v", messageType)
	}
	return
}

func (plugin *Plugin) SendMsg(mt rpc.MessageType, message map[string]interface{}) {
	message["pluginId"] = plugin.pluginId
	data, err := json.Marshal(message)
	if err != nil {
		return
	}
	err = plugin.clint.WriteMessage(&rpc.BaseMessage{
		MessageType: mt,
		Data:        data,
	})
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
	plugin.SendMsg(rpc.MessageType_PluginUnloadRequest, map[string]interface{}{})
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
