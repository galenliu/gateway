package mdns

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"github.com/galenliu/gateway/pkg/matter/inet/udp_endpoint"
	"github.com/galenliu/gateway/pkg/system"
	"net/netip"
	"sync"
)

type PacketDelegate interface {
	OnMdnsPacketData(data *core.BytesRange, info *inet.IPPacketInfo)
}

//type Server interface {
//	Shutdown()
//	DirectSend() error
//	BroadcastUnicastQuery(data []byte)
//	BroadcastSend([]byte)
//	ShutdownEndpoint(aEndpoint EndpointInfo)
//	IsListening() bool
//}

type Server struct {
	mResponseDelegate PacketDelegate
	mQueryDelegate    PacketDelegate
}

var insServer *Server
var serOnce sync.Once

func GlobalServer() *Server {
	serOnce.Do(func() {
		insServer = newMdnsServer()
	})
	return insServer
}

func newMdnsServer() *Server {
	return &Server{}
}

func (m Server) Shutdown() {

}

func (m *Server) StartServer(mgr udp_endpoint.UDPEndpoint, port int) error {
	m.Shutdown()
	return m.Listen(mgr, port)
}

func (m *Server) OnQuery(data *core.BytesRange, info *inet.IPPacketInfo) {
	if m.mQueryDelegate != nil {
		m.mResponseDelegate.OnMdnsPacketData(data, info)
	}
}

func (m *Server) OnUdpPacketReceived(msg *system.PacketBufferHandle, info *inet.IPPacketInfo) {}

func (m *Server) OnReceiveErrorFunct(err error, info *inet.IPPacketInfo) {}

func (m *Server) Listen(udpEndPoint udp_endpoint.UDPEndpoint, port int) error {
	m.Shutdown()

	err := udpEndPoint.Bind(netip.Addr{}, port, inet.InterfaceId{})
	if err != nil {
		return err
	}
	err = udpEndPoint.Listen(m.OnUdpPacketReceived, m.OnReceiveErrorFunct, nil)
	if err != nil {
		return err
	}
	return nil
}

func (m *Server) SetQueryDelegate(delegate PacketDelegate) {
	m.mQueryDelegate = delegate
}
