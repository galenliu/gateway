package plugin

import (
	"github.com/galenliu/gateway/pkg/rpc"
	"github.com/galenliu/gateway/plugin/internal"
	json "github.com/json-iterator/go"
)

type GetValueFunc func() (interface{}, error)
type SetValueFunc func(v interface{}) (interface{}, error)
type NotifyPropertyChanged func(value interface{})

type Property struct {
	*internal.Property

	device       Device
	response     chan interface{}
	SetValueFunc SetValueFunc
}

func NewPropertyFormString(desc string) *Property {
	return nil
}

func (p *Property) DoPropertyChanged(data []byte) {
	value := json.Get(data, "value").GetInterface()
	err := p.Property.SetValue(value)
	if err != nil {
		return
	}

}

func (p *Property) SetValue(value interface{}) (interface{}, error) {
	var data map[string]interface{}
	data["deviceId"] = p.device.ID
	data["propertyName"] = p.Name
	data["propertyValue"] = value
	p.device.adapter.SendMessage(rpc.MessageType_DeviceSetPropertyCommand, data)
	return nil, nil
}
