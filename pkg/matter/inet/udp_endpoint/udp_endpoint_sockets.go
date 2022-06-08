package udp_endpoint

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"github.com/galenliu/gateway/pkg/system"
	"net"
	"net/netip"
)

type UDPEndPointImplSockets struct {
	inet.InterfaceId
	mBoundPort   int
	mBoundIntfId inet.InterfaceId
	mConn        *net.UDPConn
}

func NewUDPEndPointImplSockets() *UDPEndPointImplSockets {
	return &UDPEndPointImplSockets{}
}

func (s *UDPEndPointImplSockets) IPv4JoinLeaveMulticastGroupImpl(aInterfaceId inet.InterfaceId, addr netip.Addr, b bool) error {
	//TODO implement me
	panic("implement me")
}

func (s *UDPEndPointImplSockets) IPv6JoinLeaveMulticastGroupImpl(aInterfaceId inet.InterfaceId, addr netip.Addr, b bool) error {
	//TODO implement me
	panic("implement me")
}

func (s *UDPEndPointImplSockets) SendMsgImpl(pktInfo *inet.IPPacketInfo, msg *system.PacketBufferHandle) error {
	//TODO implement me
	panic("implement me")
}

func (s *UDPEndPointImplSockets) CloseImpl() {
	//TODO implement me
	panic("implement me")
}

func (s *UDPEndPointImplSockets) BindImpl(addr netip.Addr, port int, interfaceId inet.InterfaceId) error {
	if addr.Is6() {
		err := s.ipV6Bind(addr, port, interfaceId)
		if err != nil {
			return err
		}
	} else if addr.Is4() {
		err := s.ipV4Bind(addr, port, interfaceId)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("wrong address type")
	}

	s.mBoundPort = port
	s.mBoundIntfId = interfaceId
	return nil
}

func (s *UDPEndPointImplSockets) ListenImpl() error {
	if s.mConn == nil {
		return fmt.Errorf("conn err")
	}
	return nil
}

func (s *UDPEndPointImplSockets) ipV6Bind(addr netip.Addr, port int, id inet.InterfaceId) error {
	udpAddr := net.UDPAddrFromAddrPort(netip.AddrPortFrom(addr, uint16(port)))
	conn, err := net.ListenMulticastUDP("udp", &id.Interface, udpAddr)
	if err != nil {
		return err
	}
	s.mConn = conn
	return nil
}

func (s *UDPEndPointImplSockets) ipV4Bind(addr netip.Addr, port int, id inet.InterfaceId) error {
	udpAddr := net.UDPAddrFromAddrPort(netip.AddrPortFrom(addr, uint16(port)))
	conn, err := net.ListenMulticastUDP("udp", &id.Interface, udpAddr)
	if err != nil {
		return err
	}
	s.mConn = conn
	return nil
}

func (s *UDPEndPointImplSockets) getSocket() error {
	return nil
}
