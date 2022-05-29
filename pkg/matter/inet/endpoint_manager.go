package inet

import (
	"net/netip"
)

type OnMessageReceivedFunct func([]byte)
type OnReceiveErrorFunct func() []byte

type EndpointManager interface {
}

type UDPEndpointManager struct {
}

func (m UDPEndpointManager) Bind(addr netip.Addr, port int) {

}

func (m UDPEndpointManager) Listen(onMessageReceived OnMessageReceivedFunct, onReceiveError OnReceiveErrorFunct) {

}
