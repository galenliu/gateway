package app

import "github.com/galenliu/gateway/pkg/matter/server"

type Config struct {
	ConfigNetworkLayerBle bool
	mSecuredServicePort   int
	mUnsecuredServicePort int
}

func AppMainInit(con Config) {

}

func AppMainLoop() {
	s := server.NewCHIPServer()
	s.Init(server.Config{})
}
