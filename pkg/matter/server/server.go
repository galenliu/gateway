package server

import "github.com/galenliu/gateway/pkg/dnssd"

type Config struct {
	ChipDeviceConfigEnableDnssd bool
}

type CHIPServer struct {
	mSecuredServicePort   int
	mUnsecuredServicePort int
	config                Config
	dnssd                 *dnssd.Dnssd
}

func NewCHIPServer() *CHIPServer {
	return &CHIPServer{}
}

func (chip CHIPServer) Init(secureServicePort, unsecureServicePort int) {
	chip.mUnsecuredServicePort = unsecureServicePort
	chip.mSecuredServicePort = secureServicePort
	if chip.config.ChipDeviceConfigEnableDnssd {
		chip.dnssd = dnssd.NewDnssd(chip.mSecuredServicePort, chip.mUnsecuredServicePort)
	}
}
