package transport

import (
	"github.com/galenliu/gateway/pkg/inet/udp_endpoint"
	"net/netip"
)

type UdpListenParameters struct {
	mAddr            netip.Addr
	mPort            int
	mNativeParams    func()
	mEndPointManager udp_endpoint.UDPEndpoint
}

func (p *UdpListenParameters) SetListenPort(port int) {
	p.mPort = port
}

func (p *UdpListenParameters) SetNativeParams(params func()) {
	p.mNativeParams = params
}

type UdpTransport struct {
}

func NewUdpTransport(mgr udp_endpoint.UDPEndpoint, params UdpListenParameters) (*UdpTransport, error) {
	return nil, nil
}
