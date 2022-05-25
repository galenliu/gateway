package server

import (
	"github.com/galenliu/gateway/pkg/matter/credentials"
	"sync"
)

type Config struct {
	UnsecureServicePort int
	SecureServicePort   int
}

type Server struct {
	mSecuredServicePort            int
	mUnsecuredServicePort          int
	mOperationalServicePort        int
	mUserDirectedCommissioningPort int
	interfaceId                    any
	config                         Config
	dnssdServer                    *DnssdServer
	mFabrics                       *credentials.FabricTable
}

func NewCHIPServer() *Server {
	return &Server{}
}

var ins *Server
var once sync.Once

func GetInstance() *Server {
	once.Do(func() {
		ins = NewCHIPServer()
	})
	return ins
}

func (chip Server) Init(initParams ServerInitParams) {
	chip.mUnsecuredServicePort = initParams.operationalServicePort
	chip.mSecuredServicePort = initParams.userDirectedCommissioningPort

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
