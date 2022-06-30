package mdns

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/inet/IP"
	"github.com/galenliu/gateway/pkg/inet/IPPacket"
	"github.com/galenliu/gateway/pkg/inet/Interface"
	"github.com/galenliu/gateway/pkg/inet/udp_endpoint"
	"github.com/galenliu/gateway/pkg/system"
	"net/netip"
	"sync"
)

type Server struct {
	mResponseDelegate     PacketDelegate
	mQueryDelegate        PacketDelegate
	mIpv6BroadcastAddress netip.Addr
}

func (m *Server) ShutdownEndpoint(info EndpointInfo) {
	//TODO implement me
	panic("implement me")
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

func (m *Server) OnQuery(data *core.BytesRange, info *IPPacket.Info) {
	if m.mQueryDelegate != nil {
		m.mResponseDelegate.OnMdnsPacketData(data, info)
	}
}

func (m *Server) OnUdpPacketReceived(msg *system.PacketBufferHandle, info *IPPacket.Info) {}

func (m *Server) OnReceiveErrorFunct(err error, info *IPPacket.Info) {}

func (m *Server) Listen(udpEndPoint udp_endpoint.UDPEndpoint, port int) error {
	m.Shutdown()

	err := udpEndPoint.Bind(netip.Addr{}, port, Interface.Id{})
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

func (m *Server) SetDelegate() {

}

func (m *Server) DirectSend(packet *system.PacketBufferHandle, address IP.Address, port int, id Interface.Id) error {
	return nil
}

func (m *Server) BroadcastSend(packet *system.PacketBufferHandle, port int, id Interface.Id, addr IP.Address) error {
	socketPicker := ListenSocketPickerDelegate{}
	filter := InterfaceTypeFilterDelegate{
		mChild:     socketPicker,
		mInterface: id,
		mAddress:   addr.Addr,
	}

	return m.broadcastImpl(packet, port, filter)
}

func (m *Server) broadcastImpl(packet *system.PacketBufferHandle, port int, delegate InterfaceTypeFilterDelegate) error {
	//successes := 0
	//failures := 0
	//var err error
	var info EndpointInfo
	udp := delegate.Accept(info)
	if info.mAddress.Is6() {
		return udp.SendTo(m.mIpv6BroadcastAddress, port, packet, udp.GetBoundInterface())

	}
	return nil
}

func GetIpv4Into() IP.Address {
	addr := netip.AddrFrom4([4]byte{224, 0, 0, 251})
	return IP.Address{Addr: addr}
}

func GetIpv6Into() IP.Address {
	addr, _ := netip.ParseAddr("FF02::FB")
	return IP.Address{Addr: addr}
}
