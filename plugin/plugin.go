// 处理plugin的消息，完成plugin生命周期状态管理
package plugin

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway-grpc"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rpc"
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
	Clint         rpc.Clint
	logger        logging.Logger
	addonManager  *Manager
	adapters      sync.Map
	notifiers     sync.Map
	apiHandlers   sync.Map
	services      sync.Map
}

func NewPlugin(pluginId string, manager *Manager, s *PluginsServer, log logging.Logger) (plugin *Plugin) {
	plugin = &Plugin{}
	plugin.logger = log
	plugin.locker = new(sync.Mutex)
	plugin.addonManager = manager
	plugin.closeChan = make(chan struct{})
	plugin.closeExecChan = make(chan struct{})
	plugin.pluginId = pluginId
	plugin.restart = false
	plugin.pluginServer = s
	plugin.execPath = path.Join(plugin.pluginServer.manager.config.UserProfile.AddonsDir, pluginId)
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

func (plugin *Plugin) OnMsg(messageType gateway_grpc.MessageType, data []byte) (err error) {

	switch messageType {
	case gateway_grpc.MessageType_AdapterAddedNotification:
		var message gateway_grpc.AdapterAddedNotificationMessage_Data
		err = json.Unmarshal(data, &message)
		if err != nil {
			return err
		}
		adapter := NewAdapter(plugin, message.AdapterId, message.Name, message.PackageName, plugin.logger)
		plugin.adapters.Store(adapter.ID, adapter)
		plugin.addonManager.addAdapter(adapter)
		return nil

	case gateway_grpc.MessageType_NotifierAddedNotification:
		return

	}
	adapterId := json.Get(data, "adapterId").ToString()
	adapter := plugin.getAdapter(adapterId)
	if adapter == nil {
		return fmt.Errorf("(%s)adapter not found", gateway_grpc.MessageType_name[int32(gateway_grpc.MessageType_AdapterAddedNotification)])
	}
	// adapter handler
	{
		switch messageType {
		case gateway_grpc.MessageType_DeviceRequestActionResponse:
			break
		case gateway_grpc.MessageType_DeviceRemoveActionResponse:
			break
		case gateway_grpc.MessageType_OutletNotifyResponse:
			break
		case gateway_grpc.MessageType_AdapterUnloadResponse:
			break
		case gateway_grpc.MessageType_DeviceSetCredentialsResponse:
			break
		case gateway_grpc.MessageType_ApiHandlerApiResponse:
			break

		case gateway_grpc.MessageType_ApiHandlerAddedNotification:
			return
		case gateway_grpc.MessageType_ApiHandlerUnloadResponse:
			return
		case gateway_grpc.MessageType_PluginUnloadRequest:
			return
		case gateway_grpc.MessageType_PluginErrorNotification:
			return
		case gateway_grpc.MessageType_DeviceAddedNotification:
			//messages.DeviceAddedNotification
			var message gateway_grpc.DeviceAddedNotificationMessage_Data
			err = json.Unmarshal(data, &message)
			if err != nil {
				return err
			}
			var newDevice = NewDeviceFormString(string(message.Device), adapter)
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
		case gateway_grpc.MessageType_ServiceAddedNotification:
			var d gateway_grpc.ServiceAddedNotificationMessage_Data
			err := json.Unmarshal(data, &d)
			if err != nil {
				return err
			}
			newService := NewService(plugin, plugin.addonManager.Eventbus, d.ServiceId, d.Name)
			plugin.services.Store(newService.ID, newService)
			plugin.addonManager.addService(newService)

		case gateway_grpc.MessageType_ServiceGetThingsRequest:
			things := plugin.pluginServer.manager.container.GetThings()
			bt, err := json.Marshal(things)
			if err != nil {
				return err
			}
			var d = make(map[string]interface{})
			d["things"] = bt
			adapter.sendMsg(gateway_grpc.MessageType_ServiceGetThingsResponse, d)
		case gateway_grpc.MessageType_ServiceGetThingRequest:
			id := json.Get(data, "thingId").ToString()
			thing := plugin.pluginServer.manager.container.GetThing(id)
			bt, err := json.Marshal(thing)
			if err != nil {
				return err
			}
			var d = make(map[string]interface{})
			d["thing"] = bt
			adapter.sendMsg(gateway_grpc.MessageType_ServiceGetThingResponse, d)
		case gateway_grpc.MessageType_ServiceSetPropertyValueRequest:
			var message gateway_grpc.ServiceSetPropertyValueRequestMessage_Data
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
		case gateway_grpc.MessageType_AdapterUnloadResponse:
			return

		case gateway_grpc.MessageType_NotifierUnloadResponse:
			return

		case gateway_grpc.MessageType_AdapterRemoveDeviceResponse:
			device := adapter.plugin.pluginServer.manager.getDevice(json.Get(data, "deviceId").ToString())
			adapter.handleDeviceRemoved(device)

		case gateway_grpc.MessageType_OutletAddedNotification:
			return
		case gateway_grpc.MessageType_OutletRemovedNotification:
			return

		case gateway_grpc.MessageType_DeviceSetPinResponse:

		case gateway_grpc.MessageType_DevicePropertyChangedNotification:
			var message gateway_grpc.DevicePropertyChangedNotificationMessage_Data
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

		case gateway_grpc.MessageType_DeviceActionStatusNotification:
			var message gateway_grpc.DeviceActionStatusNotificationMessage_Data
			err = json.Unmarshal(data, &message)
			if err != nil {
				return err
			}
			device := adapter.plugin.pluginServer.manager.getDevice(message.DeviceId)
			device.actionNotify(&message)

		case gateway_grpc.MessageType_DeviceEventNotification:
			var event addon.Event
			json.Get(data, "event").ToVal(&event)

		case gateway_grpc.MessageType_DeviceConnectedStateNotification:
			var message gateway_grpc.DeviceConnectedStateNotificationMessage_Data
			err = json.Unmarshal(data, &message)
			if err != nil {
				return err
			}
			device := adapter.plugin.pluginServer.manager.getDevice(message.DeviceId)
			device.connectedNotify(message.Connected)

		case gateway_grpc.MessageType_AdapterPairingPromptNotification:
			return

		case gateway_grpc.MessageType_AdapterUnpairingPromptNotification:
			return
		case gateway_grpc.MessageType_MockAdapterClearStateResponse:
			return

		case gateway_grpc.MessageType_MockAdapterRemoveDeviceResponse:
			return
		default:
			return nil
		}
	}
	return nil
}

func (plugin *Plugin) SendMsg(mt gateway_grpc.MessageType, message map[string]interface{}) {
	message["pluginId"] = plugin.pluginId
	data, err := json.Marshal(message)
	if err != nil {
		return
	}
	err = plugin.Clint.Send(&gateway_grpc.BaseMessage{
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
	plugin.restart = false
	plugin.unloading = true
	plugin.SendMsg(gateway_grpc.MessageType_PluginUnloadRequest, map[string]interface{}{})
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
