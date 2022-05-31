package inet

import (
	"net"
	"net/netip"
)

type UDPEndpoint struct {
	mInterface net.Interface
	mAddr      netip.Addr
	mPort      int
}

func (e *UDPEndpoint) Bind(addr netip.Addr, port int, interfaceId net.Interface) {
	e.mAddr = addr
	e.mInterface = interfaceId
	e.mPort = port
}
