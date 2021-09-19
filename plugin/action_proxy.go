package plugin

import "github.com/galenliu/gateway-addon/actions"

type Action struct {
	device *Device
	*actions.Action
}

func NewAction(device *Device, action *actions.Action) *Action {
	return &Action{
		device: device,
		Action: action,
	}
}
