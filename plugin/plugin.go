// 处理plugin的消息，完成plugin生命周期状态管理
package plugin

import (
	"context"
	"fmt"
	addon "github.com/galenliu/gateway-addon"
	"github.com/galenliu/gateway/configs"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/plugin/internal"
	json "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

const ExecNode = "{nodeLoader}"
const ExecPython3 = "{python}"

type Plugin struct {
	locker        *sync.Mutex
	pluginId      string
	exec          string
	execPath      string
	registered    bool
	conn          *Connection
	closeChan     chan struct{}
	closeExecChan chan struct{}
	pluginServer  *PluginsServer
	logger        logging.Logger
}

func NewPlugin(s *PluginsServer, pluginId string, log logging.Logger) (plugin *Plugin) {
	plugin = &Plugin{}
	plugin.logger = log
	plugin.locker = new(sync.Mutex)
	plugin.closeChan = make(chan struct{})
	plugin.closeExecChan = make(chan struct{})
	plugin.pluginId = pluginId
	plugin.registered = false
	plugin.pluginServer = s
	plugin.execPath = path.Join(plugin.pluginServer.manager.options.DataDir, pluginId)
	return
}

//当一个plugin建立连接后，则回复网关数据。
func (plugin *Plugin) handleConnection(c *Connection, d []byte) {
	if d != nil {
		plugin.handleMessage(d)
	}
	for {
		select {
		case <-plugin.closeChan:
			data := make(map[string]interface{})
			plugin.send(internal.PluginUnloadResponse, data)
			return
		default:
			_, data, err := c.ws.ReadMessage()
			if err != nil {
				plugin.logger.Error("plugin read err :%s", err.Error())
				plugin.registered = false
				return
			}
			plugin.handleMessage(data)

		}
	}
}

//传入的data=序列化后的 Message.Data
func (plugin *Plugin) handleMessage(data []byte) {

	var messageType int
	if tp := json.Get(data, "messageType"); tp.ValueType() == json.NumberValue {
		messageType = tp.ToInt()
	} else {
		return
	}

	//如果为0，则消息不合法(如：缺少 messageType字段)
	if messageType == 0 {
		plugin.logger.Info("messageType err")
		return
	}
	var adapterId = json.Get(data, "data", "adapterId").ToString()

	// plugin handler
	switch messageType {
	case internal.DeviceRequestActionResponse:
		break
	case internal.DeviceRemoveActionResponse:
		break
	case internal.OutletNotifyResponse:
		break

	case internal.AdapterUnloadResponse:
		break
	case internal.DeviceSetCredentialsResponse:
		break
	case internal.ApiHandlerApiResponse:
		break

	}

	switch messageType {
	case internal.AdapterAddedNotification:
		var name = json.Get(data, "data", "name").ToString()
		var packageName = json.Get(data, "data", "packageName").ToString()
		if packageName == "" {
			return
		}
		adapter := NewAdapter(plugin, name, adapterId, packageName, plugin.logger)
		plugin.pluginServer.manager.addAdapter(adapter)
		return
	}

	adapter := plugin.pluginServer.manager.getAdapter(adapterId)
	if adapter == nil {
		plugin.logger.Info("(%s)adapter not found", internal.MessageTypeToString(int(messageType)))
		return
	}

	switch messageType {

	case internal.NotifierAddedNotification:
		return
	case internal.ApiHandlerAddedNotification:
		return
	case internal.ApiHandlerUnloadResponse:
		return
	case internal.PluginUnloadRequest:
		return
	case internal.PluginErrorNotification:
		return
	case internal.DeviceAddedNotification:
		//messages.DeviceAddedNotification

		data := gjson.GetBytes(data, "data").Get("device").String()
		if data == "" {
			plugin.logger.Info("marshal device err")
			return
		}
		var newDevice = internal.NewDeviceFormString(data)

		if newDevice == nil {
			plugin.logger.Error("device add err:")
			return
		}
		newDevice.AdapterId = adapterId
		adapter.handleDeviceAdded(newDevice)
		return

	}

	deviceId := json.Get(data, "data", "deviceId").ToString()
	device := adapter.getDevice(deviceId)
	if device == nil {
		plugin.logger.Info("device cannot found: %s", deviceId)
		return
	}

	switch messageType {
	case internal.AdapterUnloadResponse:
		return

	case internal.NotifierUnloadResponse:
		return

	case internal.AdapterRemoveDeviceResponse:
		adapter.handleDeviceRemoved(device)

	case internal.OutletAddedNotification:
		return
	case internal.OutletRemovedNotification:
		return

	case internal.DeviceSetPinResponse:
		s := json.Get(data, "pin").ToString()
		var pin addon.PIN
		err := json.UnmarshalFromString(s, &pin)
		if err != nil {
			plugin.logger.Info("pin error")
			return
		}
		ee := device.SetPin(pin)
		if ee != nil {
			plugin.logger.Info(ee.Error())
		}

	case internal.DevicePropertyChangedNotification:

		prop := gjson.GetBytes(data, "data.property").String()
		propName := gjson.GetBytes(data, "data.property.name").String()

		property := device.GetProperty(propName)
		if property == nil {
			logging.Info("propName err")
			return
		}
		if len(prop) == 0 {
			return
		}
		property.DoPropertyChanged(prop)
		Publish(constant.PropertyChanged, property.AsDict())
		return

	case internal.DeviceActionStatusNotification:

		json.Get(data, "data", "action").ToVal(&action)
		return

	case internal.DeviceEventNotification:

		json.Get(data, "data", "event").ToVal(&event)

	case internal.DeviceConnectedStateNotification:
		var connected = json.Get(data, "data", "connected")
		if connected.LastError() == nil {
			event_bus.Publish(constant.CONNECTED, device, connected.ToBool())
		}
		return

	case internal.AdapterPairingPromptNotification:
		return

	case internal.AdapterUnpairingPromptNotification:
		return
	case internal.MockAdapterClearStateResponse:
		return

	case internal.MockAdapterRemoveDeviceResponse:
		return

	}
}

func (plugin *Plugin) execute() {
	plugin.exec = strings.Replace(plugin.exec, "\\", string(os.PathSeparator), -1)
	plugin.exec = strings.Replace(plugin.exec, "/", string(os.PathSeparator), -1)
	command := strings.Replace(plugin.exec, "{path}", plugin.execPath, 1)
	command = strings.Replace(command, "{nodeLoader}", configs.GetNodeLoader(), 1)
	if !strings.HasPrefix(command, "python") {
		plugin.logger.Error("Now only support plugin with python lang")
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
						plugin.logger.Info(fmt.Sprintf("plugin(%s) out: %s \t\n", plugin.pluginId, string(outputByte)))
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
	//stdOut, err := cmd.StdoutPipe()

	go func() {
		err := cmd.Start()
		if err != nil {
			plugin.logger.Info("plugin(%s) run err: %s", plugin.pluginId, err.Error())
			return
		}
	}()

	plugin.logger.Debug(fmt.Sprintf("plugin(%s) execute \t\n", plugin.pluginId))
	go syncLog(stdout)
	go syncLog(stderr)

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
}

func (plugin *Plugin) unload() {
	plugin.registered = false
	plugin.closeChan <- struct{}{}
	time.Sleep(2 * time.Second)
	plugin.closeExecChan <- struct{}{}
}

func (plugin *Plugin) registerAndHandleConnection(c *Connection) {
	if plugin.registered == true {
		plugin.logger.Error("plugin is registered")
		return
	}
	plugin.conn = c
	data := make(map[string]interface{})
	data["gatewayVersion"] = configs.GetGatewayVersion()
	data["userProfile"] = configs.GetUserProfile()
	data["preferences"] = configs.GetPreferences()
	plugin.send(internal.PluginRegisterResponse, data)
	plugin.registered = true
	plugin.logger.Info("plugin: %s registered", plugin.pluginId)
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
	if configs.IsVerbose() {
		plugin.logger.Debug("Send-- %s : \t\n %s", internal.MessageTypeToString(mt), bt)
	}

	if err != nil {
		plugin.logger.Error(err.Error())
		return
	}
	plugin.conn.send(bt)
}
