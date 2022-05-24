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
	con := config.GetInstance()
	s := server.GetInstance()

	s.Init(con.SecuredDevicePort, con.SecuredCommissionerPort)
}
