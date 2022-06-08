package inet

import (
	"net/netip"
)

type IPPacketInfo struct {
	SrcAddress  netip.Addr
	DestAddress netip.Addr
	InterfaceId InterfaceId
	SrcPort     int
	DestPort    int
}
