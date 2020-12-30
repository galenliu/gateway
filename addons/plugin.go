package addons

import (
	"context"
	"fmt"
	"gateway/pkg/log"
	addon "gitee.com/liu_guilin/gateway-addon-golang"
	json "github.com/json-iterator/go"
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
	plugin.adapters = make(map[string]*AdapterProxy, 30)
	plugin.execPath = path.Join(s.manager.AddonsDir, pluginId)
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

	var adapterId = json.Get(data, "data", "adapterId").ToString()
	// plugin handler
	switch messageType {
	//adapter add notify
	case AdapterAddedNotification:
		var name = json.Get(data, "data", "name").ToString()
		var packetName = json.Get(data, "data", "packageName").ToString()
		if packetName == "" {
			return
		}
		adapter := NewAdapterProxy(plugin.pluginServer.manager, plugin, adapterId, name, packetName)
		plugin.addAdapter(adapter)
		plugin.pluginServer.addAdapter(adapter)
		log.Info(fmt.Sprintf("adapter：%s added", adapterId))
		break
	}

	//adapter handler
	adapter, ok := plugin.adapters[adapterId]
	if !ok {
		log.Error("adapter(%s) not registered", adapterId)
		return
	}

	switch messageType {
	//add deviceProxy
	case DeviceAddedNotification:
		//messages.DeviceAddedNotification
		deviceInfo := json.Get(data, "data", "device").ToString()

		var model NewDeviceModel
		err := json.UnmarshalFromString(deviceInfo, &model)
		if err != nil {
			log.Info("deviceProxy unmarshal err : %s", err.Error())
			return
		}
		device := NewDeviceProxy(adapter, &model)
		adapter.handlerDeviceAdded(device)
		break

	//deviceProxy property changed notify
	case DevicePropertyChangedNotification:
		deviceId := json.Get(data, "data", "deviceId").ToString()
		deviceProxy, ok := plugin.pluginServer.manager.devices[deviceId]
		if !ok {
			log.Info("deviceProxy(%s) not find", deviceId)
		}
		propInfo := json.Get(data, "data", "property").ToString()
		//var m messages.DevicePropertyChangedNotification
		var newProp addon.Property
		_ = json.UnmarshalFromString(propInfo, &newProp)
		propProxy, ok := deviceProxy.Properties[newProp.Name]
		if !ok {
			log.Info("deviceProxy(%s) not find prop(%s)", deviceId, propProxy.Name)
			return
		}
		propProxy.DoPropertyChanged(&newProp)
		deviceProxy.notifyPropertyChanged(propProxy)
		break
	}
}

func (plugin *Plugin) sendData(data []byte) {
	plugin.conn.send(data)
}

func (plugin *Plugin) addAdapter(adapter *AdapterProxy) {
	plugin.locker.Lock()
	defer plugin.locker.Unlock()
	plugin.adapters[adapter.ID] = adapter
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
		log.Error("plugin("+plugin.pluginId+") start failed,err: ", err.Error())
	} else {
		log.Info("start plugin:(%s)", plugin.pluginId)
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
