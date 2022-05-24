package mdns

import "net/netip"

type EndpointInfo struct {
	network    netip.Addr
	enableIPV4 bool
}

func (info EndpointInfo) IsIpv6() bool {
	return info.network.Is6()
}

type MdnsServer struct {
}

func NewMdnsServer() *MdnsServer {
	return &MdnsServer{}
}

func (m MdnsServer) Shutdown() {

}

func (m MdnsServer) StartServer() error {
	return nil
}
