package models

import (
	"gateway/pkg/thing"
	"gateway/plugin"
	"github.com/brutella/hc/service"
)

type HCServiceProxy struct {
	Service *service.Service
	DeviceID string
}

func (s *HCServiceProxy) NewHCService(typ string) {

	switch typ {
	case thing.Light:
		sev:=service.NewLightbulb()
		sev.On.OnValueRemoteUpdate(s.OnBoolValueChanged)
	case thing.OnOffSwitch:
		sev:= service.NewSwitch()
		sev.On.OnValueRemoteUpdate(s.OnBoolValueChanged)
	}

}

func (s *HCServiceProxy) OnBoolValueChanged(value bool)  {
	plugin.SetProperValue(s.DeviceID,"",value)
}

