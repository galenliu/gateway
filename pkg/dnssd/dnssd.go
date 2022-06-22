package dnssd

import (
	"log"
)

type DnssdServer struct {
	securedServicePort   int
	unsecuredServicePort int
}

// NewDnssd Dnssd初始化
func NewDnssd(mSecuredServicePort, mUnsecuredServicePort int) *DnssdServer {
	return &DnssdServer{
		securedServicePort:   mSecuredServicePort,
		unsecuredServicePort: mUnsecuredServicePort,
	}
}

// StartServer 开启Dnssd服务
func (d DnssdServer) StartServer() {
	log.Printf("")
}
