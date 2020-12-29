package addons

import (
	"gateway/pkg/log"
	json "github.com/json-iterator/go"
)

type PropertyProxy struct {
	device *DeviceProxy
	*Property
}

func (proxy *PropertyProxy) getName() string {
	return proxy.Name
}

func (proxy *PropertyProxy) setValue(value interface{}) {
	proxy.device.adapter.sendMessage(DeviceSetPropertyCommand, struct {
		AdapterId     string `json:"adapterId"`
		DeviceId      string `json:"deviceId"`
		PropertyName  string `json:"propertyName"`
		PropertyValue interface{}
	}{AdapterId: proxy.device.adapter.ID, DeviceId: proxy.device.ID, PropertyName: proxy.Name, PropertyValue: value})
}

func (proxy *PropertyProxy) AsDict() (d string) {
	d, e := json.MarshalToString(proxy)
	if e != nil {
		log.Error("property marshal err")
	}
	return d
}
