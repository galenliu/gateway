package addons

import (
	addon "gitee.com/liu_guilin/gateway-addon-golang"
)

type ActionProxy struct {
	*addon.Action
	deviceProxy *DeviceProxy
}

func NewActionProxy(devProxy *DeviceProxy, action *addon.Action) *ActionProxy {
	a := &ActionProxy{
		Action:      action,
		deviceProxy: devProxy,
	}
	return a
}
