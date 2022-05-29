package mdns

import (
	"github.com/galenliu/gateway/pkg/matter/inet"
	"net/netip"
	"sync"
)

type Server interface {
	Shutdown()
	Listen()
	DirectSend() error
	BroadcastUnicastQuery(data []byte)
	BroadcastSend([]byte)
	ShutdownEndpoint(aEndpoint EndpointInfo)
	IsListening() bool
	SetDelegate()
}

type MdnsServer struct {
}

var insServer *MdnsServer
var serOnce sync.Once

func GlobalMdnsServer() *MdnsServer {
	serOnce.Do(func() {
		insServer = newMdnsServer()
	})
	return insServer
}

func newMdnsServer() *MdnsServer {
	return &MdnsServer{}
}

func (m MdnsServer) Shutdown() {

}

func (m *MdnsServer) StartServer(mgr inet.UDPEndpointManager, port int) error {
	m.Shutdown()
	return m.Listen(mgr, port)
}

func (m *MdnsServer) OnUdpPacketReceived(updEndPoint *inet.UDPEndpoint, data []byte) {}

func (m *MdnsServer) Listen(udpEndPoint inet.UDPEndpointManager, port int) error {
	m.Shutdown()
	udpEndPoint.Bind(netip.Addr{}, port)
	//udpEndPoint.Listen(m.OnUdpPacketReceived, nil)
	return nil
}
