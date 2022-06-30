package udp_endpoint

import (
	"github.com/galenliu/gateway/pkg/errors"
	"github.com/galenliu/gateway/pkg/inet/IP"
	"github.com/galenliu/gateway/pkg/inet/IPPacket"
	"github.com/galenliu/gateway/pkg/inet/Interface"
	"github.com/galenliu/gateway/pkg/system"
	"net/netip"
)

type EndPointStateNetworkFramework struct {
	mAddrType IP.IPAddressType
}

type UDPEndPointImplNetworkFramework struct {
	*EndPointStateNetworkFramework
}

func (s *UDPEndPointImplNetworkFramework) BindImpl(addrType IP.IPAddressType, port int, interfaceId Interface.Id) error {

	if interfaceId.IsPresent() {
		return errors.NotImplement()
	}
	return nil
}

func (s *UDPEndPointImplNetworkFramework) IPv4JoinLeaveMulticastGroupImpl(aInterfaceId Interface.Id, addr netip.Addr, b bool) error {
	//TODO implement me
	panic("implement me")
}

func (s *UDPEndPointImplNetworkFramework) IPv6JoinLeaveMulticastGroupImpl(aInterfaceId Interface.Id, addr netip.Addr, b bool) error {
	//TODO implement me
	panic("implement me")
}

func (s *UDPEndPointImplNetworkFramework) SendMsgImpl(pktInfo *IPPacket.Info, msg *system.PacketBufferHandle) error {
	//TODO implement me
	panic("implement me")
}

func (s *UDPEndPointImplNetworkFramework) CloseImpl() {
	//TODO implement me
	panic("implement me")
}

func (s *UDPEndPointImplNetworkFramework) ListenImpl() error {
	return nil
}
