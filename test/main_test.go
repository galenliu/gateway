package test

import (
	"gateway"
	"gateway/addons"
	"testing"
)

func addonTest(t *testing.T){
	rt,_:= gateway.InitRuntime("")
	gw,_:= gateway.CreateGateway(rt)
	gw.AddonsManager = addons.NewAddonsManager(gw)
	_=gw.AddonsManager.InstallAddonFromUrl("yamaha-adapter",
		"https://github.com/tim-hellhake/yamaha-adapter/releases/download/0.2.0/yamaha-adapter-0.2.0.tgz",
		"6bf7adc14cfa1bda1a1f6625f2ea23e1ee7c376a127fbcba94391390f828f6d2",true)
}
