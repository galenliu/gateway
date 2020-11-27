package addons

import (
	ipc "gitee.com/liu_guilin/WebThings-schema"
	json "github.com/json-iterator/go"
)

type PropertyProxy struct {
	device *DeviceProxy
	*ipc.Property
}

func (proxy *PropertyProxy) getName() string {
	return proxy.Name
}

func (proxy *PropertyProxy) setValue(value interface{}) {
	var message = ipc.DeviceSetPropertyCommand{}
	message.DeviceId = proxy.device.ID
	message.PluginId = proxy.device.adapter.plugin.pluginId
	message.AdapterId = proxy.device.adapter.ID
	message.PropertyName = proxy.Name
	message.PropertyValue = value
	proxy.device.adapter.sendMessage(ipc.MessageTypeDeviceSetPropertyCommand, message)
}

func (proxy *PropertyProxy) AsDict() (d string) {
	d, e := json.MarshalToString(proxy)
	if e != nil {
		log.Warn("property marshal err")
	}
	return d
}