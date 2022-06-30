package udp_endpoint

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/inet/IPPacket"
	"github.com/galenliu/gateway/pkg/inet/Interface"
	"github.com/galenliu/gateway/pkg/system"
	"net"
	"net/netip"
)

type EndPointStateSockets struct {
	mAddr netip.Addr
	conn  net.Conn
}

type UDPEndPointImplSockets struct {
	*EndPointStateSockets
	Interface.Id
	mBoundPort   int
	mBoundIntfId Interface.Id
	mConn        *net.UDPConn
}

func NewUDPEndPointImplSockets() *UDPEndPointImplSockets {
	return &UDPEndPointImplSockets{}
}

func (s *UDPEndPointImplSockets) IPv4JoinLeaveMulticastGroupImpl(aInterfaceId Interface.Id, addr netip.Addr, b bool) error {
	//TODO implement me
	panic("implement me")
}

func (s *UDPEndPointImplSockets) IPv6JoinLeaveMulticastGroupImpl(aInterfaceId Interface.Id, addr netip.Addr, b bool) error {
	//TODO implement me
	panic("implement me")
}

func (s *UDPEndPointImplSockets) SendMsgImpl(pktInfo *IPPacket.Info, msg *system.PacketBufferHandle) error {
	conn, err := net.DialUDP("udp", &net.UDPAddr{
		IP:   s.mAddr.AsSlice(),
		Port: s.mBoundPort,
		Zone: "",
	}, &net.UDPAddr{
		IP:   pktInfo.DestAddress.AsSlice(),
		Port: pktInfo.DestPort,
		Zone: "",
	})
	if err != nil {
		return err
	}
	_, err = conn.Write(msg.Bytes())
	return err

}

func (s *UDPEndPointImplSockets) CloseImpl() {
	//TODO implement me
	panic("implement me")
}

func (s *UDPEndPointImplSockets) BindImpl(addr netip.Addr, port int, interfaceId Interface.Id) error {

	s.mBoundPort = port
	s.mAddr = addr
	s.mBoundIntfId = interfaceId
	return nil
}

func (s *UDPEndPointImplSockets) ListenImpl() error {
	if s.mConn == nil {
		return fmt.Errorf("conn err")
	}
	return nil
}

func (s *UDPEndPointImplSockets) ipV6Bind(addr netip.Addr, port int, id Interface.Id) error {
	udpAddr := net.UDPAddrFromAddrPort(netip.AddrPortFrom(addr, uint16(port)))
	conn, err := net.ListenMulticastUDP("udp", &id.Interface, udpAddr)
	if err != nil {
		return err
	}
	s.mConn = conn
	return nil
}

func (s *UDPEndPointImplSockets) ipV4Bind(addr netip.Addr, port int, id Interface.Id) error {
	udpAddr := net.UDPAddrFromAddrPort(netip.AddrPortFrom(addr, uint16(port)))
	conn, err := net.ListenMulticastUDP("udp", &id.Interface, udpAddr)
	if err != nil {
		return err
	}
	s.mConn = conn
	return nil
}
