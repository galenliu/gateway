package transport

import (
	"github.com/galenliu/gateway/pkg/matter/inet"
	"net/netip"
)

type UdpListenParameters struct {
	mAddr            netip.Addr
	mPort            int
	mNativeParams    func()
	mEndPointManager inet.EndpointManager
}

func (p *UdpListenParameters) SetListenPort(port int) {
	p.mPort = port
}

func (p *UdpListenParameters) SetNativeParams(params func()) {
	p.mNativeParams = params
}

type UdpTransport struct {
}

func NewUdpTransport(mgr inet.EndpointManager, params UdpListenParameters) (*UdpTransport, error) {
	return nil, nil
}
