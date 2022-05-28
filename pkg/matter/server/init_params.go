package server

import (
	"github.com/galenliu/gateway/pkg/matter/access"
	"github.com/galenliu/gateway/pkg/matter/lib"
	"net"
)

type InitParams struct {
	OperationalServicePort        int
	UserDirectedCommissioningPort int
	InterfaceId                   net.Interface
	AppDelegate                   any //unknown
	PersistentStorageDelegate     lib.PersistentStorageDelegate
	SessionResumptionStorage      any
	AccessDelegate                access.Delegate
	AclStorage                    AclStorage
	EndpointNativeParams          func()
}

func DefaultServerInitParams() *InitParams {
	return &InitParams{
		OperationalServicePort:        ChipPort,
		UserDirectedCommissioningPort: ChipUdcPort,
	}
}

func InitializeStaticResourcesBeforeServerInit() (initParams InitParams) {
	initParams = InitParams{
		OperationalServicePort:        0,
		UserDirectedCommissioningPort: 0,
	}
	list, _ := net.Interfaces()
	for _, inter := range list {
		adders, _ := inter.Addrs()
		if len(adders) > 1 {
			initParams.InterfaceId = inter
		}
	}
	return
}
