package mdns

import (
	"github.com/galenliu/gateway/pkg/inet/udp_endpoint"
	"net/netip"
)

type EndpointInfo struct {
	mInterfaceId Interface.Id
	mAddress     netip.Addr
	mListenUdp   *udp_endpoint.UDPEndpoint
}
