package addons

import (
	"context"
	"fmt"
	"gateway/pkg/log"
	messages "gitee.com/liu_guilin/WebThings-schema"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
	"sync"
)

const ExecNode = "{nodeLoader}"
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
	conn         *Connection
	pluginServer *PluginsServer
	adapters     map[string]*AdapterProxy
	ctx          context.Context
	cancelFunc   context.CancelFunc
}

func NewPlugin(s *PluginsServer, pluginId string, exec string, _ctx context.Context) (plugin *Plugin) {
	plugin = &Plugin{}
	plugin.locker = new(sync.Mutex)
	plugin.exec = exec
	plugin.pluginId = pluginId
	plugin.pluginServer = s
	plugin.adapters = make(map[string]*AdapterProxy)
	plugin.execPath = path.Join(s.addonManager.AddonsDir, pluginId)
	if _ctx != nil {
		plugin.ctx = _ctx
	} else {
		plugin.ctx = context.Background()
	}
	return
}

//传入的data=序列化后的 BaseMessage.Data
func (plugin *Plugin) OnMessage(data []byte) {

	var messageType = json.Get(data, "messageType").ToInt()
	log.Debug(fmt.Sprintf("rev message:%s", string(data)))

	// plugin消息, added
	switch messageType {
	case AdapterAddedNotification:
		var adapterId = json.Get(data, "data", "adapterId").ToString()
		var name = json.Get(data, "data", "name").ToString()
		var packetName = json.Get(data, "data", "packetName").ToString()
		adapter := NewAdapterProxy(plugin.pluginServer.addonManager,plugin, adapterId, name, packetName)
		plugin.addAdapter(adapter)
		log.Info(fmt.Sprintf("adapter：%s added",adapterId))
	}

	//再处理 adapter消息
	adapter := plugin.getAdapter(json.Get(data, "pluginId").ToString())
	switch messageType {
	//add device
	case DeviceAddedNotification:
		var m messages.DeviceAddedNotification
		_ = json.Unmarshal(data, &m)
		device := NewDeviceProxy(adapter, m.Device)
		adapter.handlerDeviceAdded(device)
		break

	//device change property
	case DevicePropertyChangedNotification:
		var m messages.DevicePropertyChangedNotification
		_ = json.Unmarshal(data, &m)
		adapter.getDevice(m.DeviceId).doPropertyChanged(m.Property)
		break
	}
}

func (plugin *Plugin) sendData(data []byte) {

	plugin.conn.send(data)
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
	commandStr := strings.Replace(commandPath, "{nodeLoader}", "node", 1)

	commands := strings.Split(commandStr, " ")

	var args []string
	if len(commands) > 1 {
		for i, arg := range commands {
			if i != 0 {
				args = append(args, arg)
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

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Print(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Error("addons start failed,err: ", zap.Error(err))
	} else {
		log.Info("start plugin: %v", plugin.pluginId)
		out, err := ioutil.ReadAll(stdOut)
		if err != nil {
			fmt.Print(err)
			msg := string(out)
			log.Debug(msg)
		}
	}
}

func (plugin *Plugin) Stop() {
	//TODO Send stop nf
	//plugin.sendData(12,"")
	plugin.cancelFunc()
}
