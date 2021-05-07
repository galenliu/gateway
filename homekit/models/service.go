package models

import (
	"github.com/brutella/hc/service"
	"github.com/galenliu/gateway/homekit"
	"github.com/galenliu/gateway/plugin"
)

type HCServiceProxy struct {
	Service  *service.Service
	DeviceID string
}

func (s *HCServiceProxy) NewHCService(typ string) {

	switch typ {
	case homekit.Light:
		sev := service.NewLightbulb()
		sev.On.OnValueRemoteUpdate(s.OnBoolValueChanged)
	case homekit.Switch:
		sev := service.NewSwitch()
		sev.On.OnValueRemoteUpdate(s.OnBoolValueChanged)
	}

}

func (s *HCServiceProxy) OnBoolValueChanged(value bool) {
	_, _ = plugin.SetProperty(s.DeviceID, "", value)
}
