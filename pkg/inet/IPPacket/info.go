package IPPacket

import (
	"github.com/galenliu/gateway/pkg/inet/IP"
	"github.com/galenliu/gateway/pkg/inet/Interface"
)

type Info struct {
	SrcAddress  IP.Address
	DestAddress IP.Address
	InterfaceId Interface.Id
	SrcPort     int
	DestPort    int
}
