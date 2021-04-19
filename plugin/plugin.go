// 处理plugin的消息，完成plugin生命周期状态管理
package plugin

import (
	"addon"
	"context"
	"fmt"
	"github.com/galenliu/gateway/config"
	"github.com/galenliu/gateway/pkg/bus"
	"github.com/galenliu/gateway/pkg/log"
	"github.com/galenliu/gateway/pkg/util"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

const ExecNode = "{nodeLoader}"
const ExecPython3 = "{python}"

type OnConnect = func(device addon.Device, bool2 bool)

type Plugin struct {
	locker   *sync.Mutex
	pluginId string
	exec     string
	execPath string

	registered   bool
	conn         *Connection
	closeChan    chan struct{}
	pluginServer *PluginsServer
}

func NewPlugin(s *PluginsServer, pluginId string) (plugin *Plugin) {
	plugin = &Plugin{}
	plugin.locker = new(sync.Mutex)
	plugin.closeChan = make(chan struct{})
	plugin.pluginId = pluginId
	plugin.registered = false
	plugin.pluginServer = s

	plugin.execPath = path.Join(plugin.getManager().AddonsDir, pluginId)
	return
}

//传入的data=序列化后的 Message.Data
func (plugin *Plugin) handleMessage(data []byte) {

	var messageType = gjson.GetBytes(data, "messageType").Uint()
	if config.IsVerbose() {
		log.Info("Read messageType: %s \t\n Data: %s", MessageTypeToString(int(messageType)), gjson.GetBytes(data, "data").String())
	}
	//如果为0，则消息不合法(如：缺少 messageType字段)
	if messageType == 0 {
		log.Info("messageType err")
		return
	}
	var adapterId = json.Get(data, "data", "adapterId").ToString()

	// plugin handler
	switch messageType {
	case DeviceRequestActionResponse:
		break
	case DeviceRemoveActionResponse:
		break
	case OutletNotifyResponse:
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
		adapter := NewAdapter(plugin.getManager(), name, adapterId, plugin.pluginId, packageName)
		adapter.plugin = plugin
		plugin.pluginServer.addAdapter(adapter)
		return
	}

	adapter := plugin.getManager().getAdapter(adapterId)
	if adapter == nil {
		log.Info("(%s)adapter not found", MessageTypeToString(int(messageType)))
		return
	}

	switch messageType {

	case NotifierAddedNotification:
		return
	case ApiHandlerAddedNotification:
		return
	case ApiHandlerUnloadResponse:
		return
	case PluginUnloadRequest:
		return
	case PluginErrorNotification:
		return
	case DeviceAddedNotification:
		//messages.DeviceAddedNotification

		data := gjson.GetBytes(data, "data").Get("device").String()
		if data == "" {
			log.Info("marshal device err")
			return
		}
		var newDevice = addon.NewDeivceFormString(data)

		if newDevice == nil {
			log.Error("device add err:")
			return
		}
		newDevice.AdapterId = adapterId
		adapter.handleDeviceAdded(newDevice)
		return

	}

	deviceId := json.Get(data, "data", "deviceId").ToString()
	device := adapter.getDevice(deviceId)
	if device == nil {
		log.Info("device cannot found: %s", deviceId)
		return
	}

	switch messageType {
	case AdapterUnloadResponse:
		return

	case NotifierUnloadResponse:
		return

	case AdapterRemoveDeviceResponse:
		adapter.handleDeviceRemoved(device)

	case OutletAddedNotification:
		return
	case OutletRemovedNotification:
		return

	case DeviceSetPinResponse:
		s := json.Get(data, "pin").ToString()
		var pin addon.PIN
		err := json.UnmarshalFromString(s, &pin)
		if err != nil {
			log.Info("pin error")
			return
		}
		ee := device.SetPin(pin)
		if ee != nil {
			log.Info(ee.Error())
		}

	case DevicePropertyChangedNotification:

		prop := gjson.GetBytes(data, "data.property").String()
		propName := gjson.GetBytes(data, "data.property.name").String()

		property := device.GetProperty(propName)
		if property == nil {
			log.Info("propName err")
			return
		}
		if len(prop) == 0 {
			return
		}
		property.DoPropertyChanged(prop)
		Publish(util.PropertyChanged, property.AsDict())
		return

	case DeviceActionStatusNotification:
		var action addon.Action
		json.Get(data, "data", "action").ToVal(&action)
		return

	case DeviceEventNotification:
		var event addon.Event
		json.Get(data, "data", "event").ToVal(&event)

	case DeviceConnectedStateNotification:
		var connected = json.Get(data, "data", "connected")
		if connected.LastError() == nil {
			bus.Publish(util.CONNECTED, device, connected.ToBool())
		}
		return

	case AdapterPairingPromptNotification:
		return

	case AdapterUnpairingPromptNotification:
		return
	case MockAdapterClearStateResponse:
		return

	case MockAdapterRemoveDeviceResponse:
		return

	}
}

func (plugin *Plugin) getManager() *AddonManager {
	return plugin.pluginServer.manager
}

func (plugin *Plugin) addAdapter(adapter *Adapter) {
	plugin.getManager().addAdapter(adapter)
}

func (plugin *Plugin) execute() {

	plugin.exec = strings.Replace(plugin.exec, "\\", string(os.PathSeparator), -1)
	plugin.exec = strings.Replace(plugin.exec, "/", string(os.PathSeparator), -1)
	command := strings.Replace(plugin.exec, "{path}", plugin.execPath, 1)
	command = strings.Replace(command, "{nodeLoader}", config.GetNodeLoader(), 1)
	if !strings.HasPrefix(command, "python") {
		log.Error("Now only support plugin with python lang")
		return
	}
	ctx, cancelFunc := context.WithCancel(context.Background())

	commands := strings.Split(command, " ")

	var syncLog = func(reader io.ReadCloser) {
		defer cancelFunc()
		for {
			select {
			case <-plugin.closeChan:
				cancelFunc()
				return
			default:
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
		}

	}

	//ctx, plugin.cancelFunc = context.WithCancel(plugin.ctx)
	var cmd *exec.Cmd
	var args = commands[1:]
	if len(args) > 0 {
		cmd = exec.CommandContext(ctx, commands[0], args...)
	} else {
		cmd = exec.CommandContext(ctx, commands[0])
	}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	//stdOut, err := cmd.StdoutPipe()

	err := cmd.Start()
	if err != nil {
		log.Info("plugin(%s) run err: %s", plugin.pluginId, err.Error())
		return
	}

	log.Debug(fmt.Sprintf("plugin(%s) execute \t\n", plugin.pluginId))
	go syncLog(stdout)
	go syncLog(stderr)

}

func (plugin *Plugin) unload() {
	//TODO Send stop nf
	data := make(map[string]interface{})
	plugin.registered = false
	plugin.send(PluginUnloadResponse, data)
	_ = plugin.conn.ws.Close()
	plugin.closeChan <- struct{}{}
	plugin.killExec()
}

func (plugin *Plugin) killExec() {

}

//当一个plugin建立连接后，则回复网关数据。
func (plugin *Plugin) handleConnection(c *Connection, d []byte) {
	if d != nil {
		plugin.handleMessage(d)
	}
	for {
		select {
		case <-plugin.closeChan:
			return
		default:
			_, data, err := c.ws.ReadMessage()
			if err != nil {
				log.Error("plugin read err :%s", err.Error())
				plugin.registered = false
				return
			}
			plugin.handleMessage(data)

		}
	}
}

func (plugin *Plugin) registerAndHandleConnection(c *Connection) {
	if plugin.registered == true {
		log.Error("plugin is registered")
		return
	}
	plugin.conn = c
	data := make(map[string]interface{})
	data["gatewayVersion"] = config.GetGatewayVersion()
	data["userProfile"] = config.GetUserProfile()
	data["preferences"] = config.GetPreferences()
	plugin.send(PluginRegisterResponse, data)
	plugin.registered = true
	log.Info("plugin: %s registered", plugin.pluginId)
	go plugin.handleConnection(c, nil)
}

func (plugin *Plugin) send(mt int, data map[string]interface{}) {
	data["pluginId"] = plugin.pluginId
	message := struct {
		MessageType int         `json:"messageType"`
		Data        interface{} `json:"data"`
	}{
		MessageType: mt,
		Data:        data,
	}
	bt, err := json.MarshalIndent(message, "", " ")
	if config.IsVerbose() {
		log.Debug("Send-- %s : \t\n %s", MessageTypeToString(mt), bt)
	}

	if err != nil {
		log.Error(err.Error())
		return
	}
	plugin.conn.send(bt)
}
