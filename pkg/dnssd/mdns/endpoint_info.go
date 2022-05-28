package mdns

import "net"

type EndpointInfo struct {
	udp *UDPEndPoint
}

type UDPEndPoint interface {
	Bind(listener net.TCPListener)
}
