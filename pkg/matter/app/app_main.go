package app

import (
	"github.com/galenliu/gateway/pkg/matter/config"
	"github.com/galenliu/gateway/pkg/matter/server"
)

type Config struct {
	ConfigNetworkLayerBle bool
	mSecuredServicePort   int
	mUnsecuredServicePort int
}

func AppMainInit(con Config) {

}

func AppMainLoop() {

	initParams := server.InitializeStaticResourcesBeforeServerInit()
	initParams.OperationalServicePort = config.ChipPort
	initParams.UserDirectedCommissioningPort = config.ChipUdcPort

	server.GetInstance().Init(initParams)

	server.GetInstance().Shutdown()

}
