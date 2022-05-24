package server

import "github.com/galenliu/gateway/pkg/matter/credentials"

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

func (chip CHIPServer) Init(con Config) {
	chip.mUnsecuredServicePort = con.UnsecureServicePort
	chip.mSecuredServicePort = con.SecureServicePort

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
