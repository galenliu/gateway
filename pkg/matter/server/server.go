package server

import (
	"github.com/galenliu/gateway/pkg/matter/credentials"
	"sync"
)

type Config struct {
	UnsecureServicePort int
	SecureServicePort   int
}

type CHIPServer struct {
	mSecuredServicePort            int
	mUnsecuredServicePort          int
	mOperationalServicePort        int
	mUserDirectedCommissioningPort int
	interfaceId                    any
	config                         Config
	dnssdServer                    *DnssdServer
	mFabrics                       *credentials.FabricTable
}

func NewCHIPServer() *CHIPServer {
	return &CHIPServer{}
}

var ins *CHIPServer
var once sync.Once

func GetInstance() *CHIPServer {
	once.Do(func() {
		ins = NewCHIPServer()
	})
	return ins
}

func (chip CHIPServer) Init(unsecureServicePort, secureServicePort int) {
	chip.mUnsecuredServicePort = unsecureServicePort
	chip.mSecuredServicePort = secureServicePort

	chip.dnssdServer = NewDnssdServer()
	chip.dnssdServer.SetFabricTable(chip.mFabrics)

	chip.dnssdServer.SetSecuredPort(chip.mOperationalServicePort)
	chip.dnssdServer.SetUnsecuredPort(chip.mUserDirectedCommissioningPort)
	chip.dnssdServer.SetInterfaceId(chip.interfaceId)
	err := chip.dnssdServer.StartServer()
	if err != nil {
		return
	}

}
