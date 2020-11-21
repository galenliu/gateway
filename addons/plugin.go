package addons

import (
	"context"
	messages "gitee.com/liu_guilin/WebThings-schema"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
	"os/exec"
	"strings"
	"sync"
)

const ExecNode= "{nodeLoader}"
const ExecPython3 = "{python}"
//Plugin 管理Adapters
//处理每一个plugin的请求
type Plugin struct {
	looker       *sync.Mutex
	pluginId     string
	exec         string
	execPath     string
	verbose      bool
	registered   bool
	ws           *Connection
	pluginServer *PluginsServer
	adapters     map[string]*AdapterProxy

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewPlugin(pluginId string, s *PluginsServer, _ctx context.Context) (plugin *Plugin) {
	plugin = &Plugin{}
	plugin.looker = new(sync.Mutex)
	plugin.pluginId = pluginId
	plugin.pluginServer = s
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
	plugin.looker.Lock()
	defer plugin.looker.Unlock()
	plugin.adapters[a.ID] = a

}

func (plugin *Plugin) getAdapter(id string) *AdapterProxy {
	plugin.looker.Lock()
	defer plugin.looker.Unlock()
	a := plugin.adapters[id]
	return a

}

func (plugin *Plugin) start() {

	command := strings.Replace(plugin.exec, "{path}", plugin.execPath, 1)

	var ctx context.Context
	ctx, plugin.cancelFunc = context.WithCancel(plugin.ctx)
	cmd := exec.CommandContext(ctx, command)
	err := cmd.Start()
	if err != nil {
		log.Warn("addons start failed")
	} else {
		log.Info("addons is running")
	}
}

func (plugin *Plugin) kill() {
	plugin.cancelFunc()
	plugin.registered = false
}
