package server

type Config struct {
	ChipDeviceConfigEnableDnssd bool
}

type CHIPServer struct {
	mSecuredServicePort   int
	mUnsecuredServicePort int
	config                Config
}

func NewCHIPServer() *CHIPServer {
	return &CHIPServer{}
}

func (chip CHIPServer) Init(secureServicePort, unsecureServicePort int) {
	chip.mUnsecuredServicePort = unsecureServicePort
	chip.mSecuredServicePort = secureServicePort
}
