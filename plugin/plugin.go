package plugin

import (
	"addon"
	"context"
	"fmt"
	"gateway/config"
	"gateway/pkg/bus"
	"gateway/pkg/log"
	"gateway/pkg/util"
	json "github.com/json-iterator/go"
	"io"
	"os/exec"
	"path"
	"strings"
	"sync"
)

const ExecNode = "{nodeLoader}"
const ExecPython3 = "{python}"

type OnConnect = func(device addon.Device, bool2 bool)

//Plugin 管理Adapters
//处理每一个plugin的请求
type Plugin struct {
	locker       *sync.Mutex
	pluginId     string
	exec         string
	execPath     string
	verbose      bool
	registered   bool
	conn         *Connection
	pluginServer *PluginsServer
	ctx          context.Context
	cancelFunc   context.CancelFunc
}

func NewPlugin(s *PluginsServer, pluginId string) (plugin *Plugin) {
	plugin = &Plugin{}
	plugin.locker = new(sync.Mutex)
	plugin.pluginId = pluginId
	plugin.pluginServer = s
	plugin.execPath = path.Join(plugin.getManager().AddonsDir, pluginId)
	plugin.ctx = context.Background()

	return
}

//传入的data=序列化后的 Message.Data
func (plugin *Plugin) handleMessage(data []byte) {

	var messageType = json.Get(data, "messageType").ToInt()
	var adapterId = json.Get(data, "data", "adapterId").ToString()
	// plugin handler
	switch messageType {
	case DeviceRequestActionResponse:
		break
	case DeviceRemoveActionResponse:
		break
	case OutletNotifyResponse:
		break
	case DeviceSetPinRequest:
		break

	case AdapterUnloadResponse:
		break
	case DeviceSetCredentialsResponse:
		break
	case ApiHandlerApiResponse:
		break

	}

	switch messageType {
	case AdapterAddedNotification:
		var name = json.Get(data, "data", "name").ToString()
		var packageName = json.Get(data, "data", "packageName").ToString()
		if packageName == "" {
			return
		}
		adapter := NewAdapterProxy(plugin.getManager(), name, adapterId, plugin.pluginId, packageName)
		adapter.plugin = plugin
		plugin.pluginServer.addAdapter(adapter)
		break

	case NotifierAddedNotification:
		break
	case ApiHandlerAddedNotification:
		break
	case ApiHandlerUnloadResponse:
		break
	case PluginUnloadRequest:
		break
	case PluginErrorNotification:
		break
	}

	var device *addon.Device
	var err error
	var adapter *AdapterProxy

	adapter, err = plugin.getManager().findAdapter(adapterId)
	if err != nil {
		log.Error(err.Error())
		return
	}
	deviceId := json.Get(data, "data", "deviceId").ToString()

	switch messageType {
	case AdapterUnloadResponse:
		break

	case NotifierUnloadResponse:
		break

	case DeviceAddedNotification:
		//messages.DeviceAddedNotification
		var newDevice addon.Device
		js := json.Get(data, "data", "device")
		if js.LastError() != nil {
			log.Error("new device err:", js.LastError().Error())
			return
		}
		err := json.UnmarshalFromString(js.ToString(), &newDevice)
		if err != nil {
			log.Error("new device err:", err.Error())
			return

		}
		adapter.HandleDeviceAdded(&newDevice)
		break

	case AdapterRemoveDeviceResponse:
		device, err = adapter.FindDevice(deviceId)
		adapter.HandleDeviceRemoved(device)

	case OutletAddedNotification:
		break
	case OutletRemovedNotification:
		break

	case DevicePropertyChangedNotification:
		var p addon.Property
		json.Get(data, "data", "property").ToVal(&p)
		p.DeviceId = deviceId
		device = adapter.GetDevice(deviceId)
		property := device.GetProperty(p.Name)
		property.Update(&p)
		bus.Publish(util.PropertyChanged, property)
		break

	case DeviceActionStatusNotification:
		var action addon.Action
		json.Get(data, "data", "action").ToVal(&action)
		device = adapter.GetDevice(deviceId)
		if device != nil && &action != nil {
			adapter.getManager().actionNotify(&action)
		}
		break

	case DeviceEventNotification:
		var event addon.Event
		json.Get(data, "data", "event").ToVal(&event)
		device = adapter.GetDevice(deviceId)
		if device != nil && &event != nil {
			adapter.getManager().eventNotify(&event)
		}

	case DeviceConnectedStateNotification:
		device = adapter.GetDevice(deviceId)
		var connected = json.Get(data, "data", "connected")
		if device != nil && connected.LastError() != nil {
			bus.Publish(util.CONNECTED, device, connected.ToBool())
		}

	case AdapterPairingPromptNotification:
		break

	case AdapterUnpairingPromptNotification:
		break
	case MockAdapterClearStateResponse:
		break

	case MockAdapterRemoveDeviceResponse:
		break

	}
}

func (plugin *Plugin) sendData(data []byte) {
	plugin.conn.send(data)
}

func (plugin *Plugin) getManager() *AddonManager {
	return plugin.pluginServer.manager
}

func (plugin *Plugin) addAdapter(adapter *AdapterProxy) {
	plugin.getManager().addAdapter(adapter)
}

func (plugin *Plugin) start() {

	command := strings.Replace(plugin.exec, "{path}", plugin.execPath, 1)
	command = strings.Replace(command, "{nodeLoader}", config.Conf.NodeLoader, 1)

	commands := strings.Split(command, " ")

	var args []string
	if len(commands) > 1 {
		for i, arg := range commands {
			if i != 0 {
				args = append(args, arg)
			}
		}
	}

	var syncLog = func(reader io.ReadCloser) {

		buf := make([]byte, 1024, 1024)
		for {
			strNum, err := reader.Read(buf)
			if strNum > 0 {
				outputByte := buf[:strNum]
				log.Info(fmt.Sprintf("plugin(%s) out: %s \t\n", plugin.pluginId, string(outputByte)))
			}
			if err != nil {
				//读到结尾
				if err == io.EOF || strings.Contains(err.Error(), "file already closed") {
					err = nil
				}
			}
		}
	}

	//ctx, plugin.cancelFunc = context.WithCancel(plugin.ctx)
	var cmd *exec.Cmd
	if len(args) > 0 {
		cmd = exec.Command(commands[0], args...)
	} else {
		cmd = exec.Command(commands[0])
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	//stdOut, err := cmd.StdoutPipe()

	go cmd.Start()

	log.Debug(fmt.Sprintf("plugin(%s) start \t\n", plugin.pluginId))
	go syncLog(stdout)
	go syncLog(stderr)

}

func (plugin *Plugin) Stop() {
	//TODO Send stop nf
	//plugin.sendData(12,"")
	plugin.cancelFunc()
}

func (plugin *Plugin) handleConnection(c *Connection) {
	plugin.conn = c
	m := struct {
		MessageType int         `json:"messageType"`
		Data        interface{} `json:"data"`
	}{
		MessageType: PluginRegisterResponse,
		Data: struct {
			PluginID       string              `json:"pluginId"`
			GatewayVersion string              `json:"gatewayVersion"`
			UserProfile    *config.UserProfile `json:"userProfile"`
			Preferences    *config.Preferences `json:"preferences"`
		}{
			PluginID:       plugin.pluginId,
			UserProfile:    config.GetUserProfile(),
			Preferences:    config.GetPreferences(),
			GatewayVersion: config.Conf.GatewayVersion,
		},
	}
	data, _ := json.MarshalIndent(m, "", " ")
	plugin.sendData(data)
	plugin.registered = true
	log.Info(fmt.Sprintf("plugin: %s registered", plugin.pluginId))
}
