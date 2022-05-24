package dnssd

import "github.com/google/martian/log"

type Dnssd struct {
	securedServicePort   int
	unsecuredServicePort int
}

// NewDnssd Dnssd初始化
func NewDnssd(mSecuredServicePort, mUnsecuredServicePort int) *Dnssd {
	return &Dnssd{
		securedServicePort:   mSecuredServicePort,
		unsecuredServicePort: mUnsecuredServicePort,
	}
}

// StartServer 开启Dnssd服务
func (d Dnssd) StartServer() {
	log.Infof("DNS-SD StartServer")
}
