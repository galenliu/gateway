package mdns

import "github.com/galenliu/gateway/pkg/matter/inet"

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

type InetLayer interface {
	NewUDPEndPoint()
}

type MdnsServer struct {
}

func NewMdnsServer() *MdnsServer {
	return &MdnsServer{}
}

func (m MdnsServer) Shutdown() {

}

func (m *MdnsServer) StartServer(mgr inet.UDPEndpointManager, port int) error {
	m.Shutdown()
	return m.Listen(mgr, port)
}

func (m *MdnsServer) OnUdpPacketReceived(data []byte) {}

func (m *MdnsServer) Listen(udpEndPoint inet.UDPEndpointManager, port int) error {
	m.Shutdown()
	return nil
}
