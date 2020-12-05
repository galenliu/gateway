package addons

import (
	"context"
	"fmt"
	messages "gitee.com/liu_guilin/WebThings-schema"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
	"sync"
)

const ExecNode= "{nodeLoader}"
const ExecPython3 = "{python}"
//Plugin 管理Adapters
//处理每一个plugin的请求
type Plugin struct {
	locker       *sync.Mutex
	pluginId     string
	exec         string
	execPath     string
	verbose      bool
	registered   bool
	ws           *Connection
	pluginServer *PluginsServer
	adapters     map[string]*AdapterProxy
	ctx          context.Context
	cancelFunc context.CancelFunc
}

func NewPlugin(s *PluginsServer,pluginId string, exec string, _ctx context.Context) (plugin *Plugin) {
	plugin = &Plugin{}
	plugin.locker = new(sync.Mutex)
	plugin.exec = exec
	plugin.pluginId = pluginId
	plugin.pluginServer = s
	plugin.execPath = path.Join(s.addonManager.AddonsDir,pluginId)
	if _ctx != nil {
		plugin.ctx = _ctx
	} else {
		plugin.ctx = context.Background()
	}
	return
}

//传入的data=序列化后的 BaseMessage.Data
func (plugin *Plugin) OnMessage(message messages.BaseMessage, data []byte) {

	log.Info("addons rev message", zap.Any("message", message))

	// plugin消息, added
	switch message.MessageType {
	case messages.MessageTypeAdepterAddedNotification:
		var m messages.AdapterAddedNotification
		_ = json.Unmarshal(data, &m)
		adapter := NewAdapterProxy(plugin.pluginServer.addonManager, m.AdapterId, m.Name, m.PackageName)
		plugin.addAdapter(adapter)

	}

	//再处理 adapter消息
	adapter := plugin.getAdapter(json.Get(data, "plugin_id").ToString())
	switch message.MessageType {
	//add device
	case messages.MessageTypeDeviceAddedNotification:
		var m messages.DeviceAddedNotification
		_ = json.Unmarshal(data, &m)
		device := NewDeviceProxy(adapter, m.Device)
		adapter.handlerDeviceAdded(device)
		break

	//device change property
	case messages.MessageTypeDevicePropertyChangedNotification:
		var m messages.DevicePropertyChangedNotification
		_ = json.Unmarshal(data, &m)
		adapter.getDevice(m.DeviceId).doPropertyChanged(m.Property)
		break
	}
}

func (plugin *Plugin) sendMessage(message messages.BaseMessage) {
	plugin.ws.SendMessage(message)
}

func (plugin *Plugin) addAdapter(a *AdapterProxy) {
	plugin.locker.Lock()
	defer plugin.locker.Unlock()
	plugin.adapters[a.ID] = a

}

func (plugin *Plugin) getAdapter(id string) *AdapterProxy {
	plugin.locker.Lock()
	defer plugin.locker.Unlock()
	a := plugin.adapters[id]
	return a

}

func (plugin *Plugin) start() {

	commandPath := strings.Replace(plugin.exec, "{path}", plugin.execPath, 1)
	commandStr := strings.Replace(commandPath, "{nodeLoader}","node", 1)

	commands :=strings.Split(commandStr," ")

	var args []string
	if len(commands)>1{
		for i,arg :=range commands{
			if i !=0{
				args = append(args, arg)
			}
		}
	}
	//ctx, plugin.cancelFunc = context.WithCancel(plugin.ctx)
	var cmd *exec.Cmd
	if len(args)>0{
		cmd = exec.Command(commands[0], args...)
	}else {
		cmd = exec.Command(commands[0])
	}


	stdOut, err :=cmd.StdoutPipe();if err != nil{
		fmt.Print(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Warn("addons start failed,err: ",zap.Error(err))
	} else {
		fmt.Printf("start plugin: %v",plugin.pluginId)
			out,err:= ioutil.ReadAll(stdOut);if err !=nil{
				fmt.Print(err)
			msg := string(out)
			fmt.Print(msg)

		}
	}
}

func (plugin *Plugin) Stop() {
	//TODO Send stop nf
	//plugin.sendMessage(12,"")
	plugin.cancelFunc()
}
