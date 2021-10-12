// 处理plugin的消息，完成plugin生命周期状态管理
package plugin

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/server"
	"github.com/galenliu/gateway/pkg/util"
	"github.com/galenliu/gateway/plugin/addon"
	json "github.com/json-iterator/go"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

type Plugin struct {
	locker        *sync.Mutex
	pluginId      string
	exec          string
	execPath      string
	restart       bool
	unloading     bool
	closeChan     chan struct{}
	closeExecChan chan struct{}
	pluginServer  *PluginsServer
	Clint         server.Clint
	logger        logging.Logger
	addonManager  *Manager
	adapters      sync.Map
	notifiers     sync.Map
	apiHandlers   sync.Map
	services      sync.Map
}

func NewPlugin(pluginId string, manager *Manager, s *PluginsServer, log logging.Logger) (plugin *Plugin) {
	execPath := s.manager.findAddon(pluginId)
	if execPath == "" {
		return nil
	}
	plugin = &Plugin{}
	plugin.logger = log
	plugin.locker = new(sync.Mutex)
	plugin.addonManager = manager
	plugin.closeChan = make(chan struct{})
	plugin.closeExecChan = make(chan struct{})
	plugin.pluginId = pluginId
	plugin.restart = false
	plugin.pluginServer = s
	plugin.execPath = execPath
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

func (plugin *Plugin) OnMsg(messageType rpc.MessageType, data []byte) (err error) {

	switch messageType {
	case rpc.MessageType_AdapterAddedNotification:
		var message rpc.AdapterAddedNotificationMessage_Data
		err = json.Unmarshal(data, &message)
		if err != nil {
			return err
		}
		adapter := NewAdapter(plugin, message.AdapterId, message.Name, message.PackageName, plugin.logger)
		plugin.adapters.Store(adapter.ID, adapter)
		plugin.addonManager.addAdapter(adapter)
		return nil

	case rpc.MessageType_NotifierAddedNotification:
		return

	}
	adapterId := json.Get(data, "adapterId").ToString()
	adapter := plugin.getAdapter(adapterId)
	if adapter == nil {
		return fmt.Errorf("(%s)adapter not found", rpc.MessageType_name[int32(rpc.MessageType_AdapterAddedNotification)])
	}
	// adapter handler
	{
		switch messageType {
		case rpc.MessageType_DeviceRequestActionResponse:
			break
		case rpc.MessageType_DeviceRemoveActionResponse:
			break
		case rpc.MessageType_OutletNotifyResponse:
			break
		case rpc.MessageType_AdapterUnloadResponse:
			break
		case rpc.MessageType_DeviceSetCredentialsResponse:
			break
		case rpc.MessageType_ApiHandlerApiResponse:
			break

		case rpc.MessageType_ApiHandlerAddedNotification:
			return
		case rpc.MessageType_ApiHandlerUnloadResponse:
			return
		case rpc.MessageType_PluginUnloadRequest:
			return
		case rpc.MessageType_PluginErrorNotification:
			return
		case rpc.MessageType_DeviceAddedNotification:
			//messages.DeviceAddedNotification
			var message rpc.DeviceAddedNotificationMessage_Data
			err = json.Unmarshal(data, &message)
			if err != nil {
				return err
			}
			var newDevice = NewDeviceFormMessage(message.Device, adapter)
			if newDevice == nil {
				return fmt.Errorf("device add failed")
			}
			adapter.handleDeviceAdded(newDevice)
			return nil
		}
	}

	// services message
	{
		switch messageType {
		case rpc.MessageType_ServiceAddedNotification:
			var d rpc.ServiceAddedNotificationMessage_Data
			err := json.Unmarshal(data, &d)
			if err != nil {
				return err
			}
			newService := NewService(plugin, plugin.addonManager.Eventbus, d.ServiceId, d.Name)
			plugin.services.Store(newService.ID, newService)
			plugin.addonManager.addService(newService)

		case rpc.MessageType_ServiceGetThingsRequest:
			things := plugin.pluginServer.manager.container.GetThings()
			bt, err := json.Marshal(things)
			if err != nil {
				return err
			}
			var d = make(map[string]interface{})
			d["things"] = bt
			adapter.sendMsg(rpc.MessageType_ServiceGetThingsResponse, d)

		case rpc.MessageType_ServiceGetThingRequest:
			id := json.Get(data, "thingId").ToString()
			thing := plugin.pluginServer.manager.container.GetThing(id)
			bt, err := json.Marshal(thing)
			if err != nil {
				return err
			}
			var d = make(map[string]interface{})
			d["thing"] = bt
			adapter.sendMsg(rpc.MessageType_ServiceGetThingResponse, d)

		case rpc.MessageType_ServiceSetPropertyValueRequest:
			var message rpc.ServiceSetPropertyValueRequestMessage_Data
			err = json.Unmarshal(data, &message)
			if err != nil {
				return err
			}
			_, err := adapter.plugin.pluginServer.manager.SetPropertyValue(message.ThingId, message.PropertyName, message.Value)
			if err != nil {
				return err
			}
		}
	}

	// device handler
	{

		switch messageType {
		case rpc.MessageType_AdapterUnloadResponse:
			return

		case rpc.MessageType_NotifierUnloadResponse:
			return

		case rpc.MessageType_AdapterRemoveDeviceResponse:
			device := adapter.plugin.pluginServer.manager.getDevice(json.Get(data, "deviceId").ToString())
			adapter.handleDeviceRemoved(device)

		case rpc.MessageType_OutletAddedNotification:
			return
		case rpc.MessageType_OutletRemovedNotification:
			return

		case rpc.MessageType_DeviceSetPinResponse:

		case rpc.MessageType_DevicePropertyChangedNotification:
			var message rpc.DevicePropertyChangedNotificationMessage_Data
			err = json.Unmarshal(data, &message)
			if err != nil {
				return err
			}
			device := adapter.plugin.pluginServer.manager.getDevice(message.DeviceId)
			property := device.GetProperty(message.Property.Name)
			if property == nil {
				return fmt.Errorf("property not found")
			}
			err := property.doPropertyChanged(message.Property)
			if err == nil {
				plugin.addonManager.Eventbus.PublishPropertyChanged(&message)
			}

		case rpc.MessageType_DeviceActionStatusNotification:
			var message rpc.DeviceActionStatusNotificationMessage_Data
			err = json.Unmarshal(data, &message)
			if err != nil {
				return err
			}
			device := adapter.plugin.pluginServer.manager.getDevice(message.DeviceId)
			device.actionNotify(&message)

		case rpc.MessageType_DeviceEventNotification:
			var event addon.Event
			json.Get(data, "event").ToVal(&event)

		case rpc.MessageType_DeviceConnectedStateNotification:
			var message rpc.DeviceConnectedStateNotificationMessage_Data
			err = json.Unmarshal(data, &message)
			if err != nil {
				return err
			}
			device := adapter.plugin.pluginServer.manager.getDevice(message.DeviceId)
			device.connectedNotify(message.Connected)

		case rpc.MessageType_AdapterPairingPromptNotification:
			return

		case rpc.MessageType_AdapterUnpairingPromptNotification:
			return
		case rpc.MessageType_MockAdapterClearStateResponse:
			return

		case rpc.MessageType_MockAdapterRemoveDeviceResponse:
			return
		default:
			return nil
		}
	}
	return nil
}

func (plugin *Plugin) SendMsg(mt rpc.MessageType, message map[string]interface{}) {
	message["pluginId"] = plugin.pluginId
	data, err := json.Marshal(message)
	if err != nil {
		return
	}
	err = plugin.Clint.Send(&rpc.BaseMessage{
		MessageType: mt,
		Data:        data,
	})
	if err != nil {
		plugin.logger.Infof("plugin send message err: %s", err.Error())
		return
	}
}

func (plugin *Plugin) start() {
	plugin.exec = strings.Replace(plugin.exec, "\\", string(os.PathSeparator), -1)
	plugin.exec = strings.Replace(plugin.exec, "/", string(os.PathSeparator), -1)

	plugin.execPath = strings.Replace(plugin.execPath, "\\", string(os.PathSeparator), -1)
	plugin.execPath = strings.Replace(plugin.execPath, "/", string(os.PathSeparator), -1)

	command := strings.Replace(plugin.exec, "{path}", plugin.execPath, 1)
	//command = strings.Replace(command, "{nodeLoader}", configs.GetNodeLoader(), 1)
	if !strings.HasPrefix(command, "python") {
		plugin.logger.Error("now only support plugin with python lang")
		return
	}
	ctx, cancelFunc := context.WithCancel(context.Background())

	commands := strings.Split(command, " ")

	var syncLog = func(reader io.ReadCloser) {
		for {
			select {
			case <-plugin.closeChan:
				return
			default:
				buf := make([]byte, 1024)
				for {
					strNum, err := reader.Read(buf)
					if strNum > 0 {
						outputByte := buf[:strNum]
						plugin.logger.Infof("%s out:[ %s ]", plugin.pluginId, string(outputByte))
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

	go func() {
		err := cmd.Start()
		if err != nil {
			plugin.logger.Infof("%s run err: %s", plugin.pluginId, err.Error())
			return
		}
	}()

	plugin.logger.Infof("%s start", plugin.pluginId)

	closeExec := func() {
		for {
			select {
			case <-plugin.closeExecChan:
				ctx.Done()
				cancelFunc()
			}
		}
	}
	go closeExec()

	//stdOut, err := cmd.StdoutPipe()
	go syncLog(stdout)
	go syncLog(stderr)
}

func (plugin *Plugin) unload() {
	plugin.disable()
	err := util.RemoveDir(plugin.execPath)
	if err != nil {
		plugin.logger.Errorf("delete plugin dir failed err:", err.Error())
	}
	err = util.RemoveDir(path.Join(plugin.pluginServer.manager.config.UserProfile.DataDir, plugin.pluginId))
	if err != nil {
		plugin.logger.Errorf("delete plugin dir failed err:", err.Error())
	}
}

func (plugin *Plugin) disable() {
	plugin.restart = false
	plugin.unloading = true
	plugin.SendMsg(rpc.MessageType_PluginUnloadRequest, map[string]interface{}{})
	plugin.unloadComponents()
	plugin.kill()
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
}

func (plugin *Plugin) kill() {
	select {
	case plugin.closeChan <- struct{}{}:
	}
	select {
	case plugin.closeExecChan <- struct{}{}:
	}
}
