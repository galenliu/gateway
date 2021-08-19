// 处理plugin的消息，完成plugin生命周期状态管理
package plugin

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway/pkg/logging"
	"github.com/galenliu/gateway/pkg/rpc"
	"github.com/galenliu/gateway/plugin/internal"
	json "github.com/json-iterator/go"
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

type IClint interface {
	Send(message *rpc.BaseMessage) error
	Read() (*rpc.BaseMessage, error)
}

type Plugin struct {
	locker        *sync.Mutex
	pluginId      string
	exec          string
	execPath      string
	registered    bool
	closeChan     chan struct{}
	closeExecChan chan struct{}
	pluginServer  *PluginsServer
	logger        logging.Logger
	Clint         IClint
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
	plugin.execPath = path.Join(plugin.pluginServer.manager.config.BaseDir, pluginId)
	return
}

func (plugin *Plugin) MessageHandler(messageType rpc.MessageType, data []byte) (err error) {
	var adapterId = json.Get(data, "adapterId").ToString()
	if adapterId == "" {
		return fmt.Errorf("pluginId is none")
	}
	// plugin handler
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

	}

	switch messageType {
	case rpc.MessageType_AdapterAddedNotification:
		var name = json.Get(data, "name").ToString()
		var packageName = json.Get(data, "packageName").ToString()
		adapter := NewAdapter(plugin, name, adapterId, packageName, plugin.logger)
		plugin.pluginServer.manager.handleAdapterAdded(adapter)
		return nil
	}

	adapter := plugin.pluginServer.manager.getAdapter(adapterId)
	if adapter == nil {
		plugin.logger.Info("(%s)adapter not found", rpc.MessageType_name[int32(rpc.MessageType_AdapterAddedNotification)])
		return nil
	}

	switch messageType {

	case rpc.MessageType_NotifierAddedNotification:
		return
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
		str := json.Get(data, "device").ToString()
		if str == "" {
			plugin.logger.Info("marshal device err")
			return
		}
		var newDevice = internal.NewDeviceFormString(str)

		if newDevice == nil {
			plugin.logger.Error("device add err:")
			return
		}
		newDevice.AdapterId = adapterId
		adapter.handleDeviceAdded(newDevice)
		return
	}

	deviceId := json.Get(data, "deviceId").ToString()
	device := adapter.plugin.pluginServer.manager.getDevice(deviceId)
	if device == nil {
		plugin.logger.Info("device cannot found: %s", deviceId)
		return
	}
	switch messageType {
	case rpc.MessageType_AdapterUnloadResponse:
		return

	case rpc.MessageType_NotifierUnloadResponse:
		return

	case rpc.MessageType_AdapterRemoveDeviceResponse:
		adapter.handleDeviceRemoved(device)

	case rpc.MessageType_OutletAddedNotification:
		return
	case rpc.MessageType_OutletRemovedNotification:
		return

	case rpc.MessageType_DeviceSetPinResponse:
		//s := json.Get(data, "pin").ToString()
		//var pin addon.PIN
		//err := json.UnmarshalFromString(s, &pin)
		//if err != nil {
		//	plugin.logger.Info("pin error")
		//	return
		//}
		//ee := device.SetPin(pin)
		//if ee != nil {
		//	plugin.logger.Info(ee.Error())
		//}

	case rpc.MessageType_DevicePropertyChangedNotification:
		prop := json.Get(data, "property").ToString()
		propName := json.Get(data, "property", "name").ToString()

		property, ok := device.Properties[propName]
		if !ok {
			plugin.logger.Info("propName err")
			return
		}
		if len(prop) == 0 {
			return
		}
		property.DoPropertyChanged([]byte(prop))
		return

	case rpc.MessageType_DeviceActionStatusNotification:
		var action internal.Action
		json.Get(data, "action").ToVal(&action)
		return

	case rpc.MessageType_DeviceEventNotification:
		var event internal.Event
		json.Get(data, "event").ToVal(&event)

	case rpc.MessageType_DeviceConnectedStateNotification:
		var connected = json.Get(data, "connected").ToBool()
		device.SetConnect(connected)
		return

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
	return nil
}

func (plugin *Plugin) execute() {
	plugin.exec = strings.Replace(plugin.exec, "\\", string(os.PathSeparator), -1)
	plugin.exec = strings.Replace(plugin.exec, "/", string(os.PathSeparator), -1)
	command := strings.Replace(plugin.exec, "{path}", plugin.execPath, 1)
	//command = strings.Replace(command, "{nodeLoader}", configs.GetNodeLoader(), 1)
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

func (plugin *Plugin) SendMessage(mt rpc.MessageType, message map[string]interface{}) {
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
		plugin.logger.Warning("plugin send message err: %s", err.Error())
		return
	}
}
