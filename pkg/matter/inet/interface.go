package inet

import "net"

type InterfaceType = uint

const (
	IUnknown InterfaceType = iota
	WiFi
	Ethernet
	Cellular
	Thread
)

type InterfaceId struct {
	net.Interface
}
