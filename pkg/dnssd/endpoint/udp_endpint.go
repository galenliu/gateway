package endpoint

import (
	"github.com/galenliu/gateway/pkg/inet/IPPacket"
	"github.com/galenliu/gateway/pkg/inet/Interface"
	"github.com/galenliu/gateway/pkg/system"
	"net/netip"
)

type OnMessageReceivedFunct = func(*system.PacketBufferHandle, *IPPacket.Info)
type OnReceiveErrorFunct = func(error, *IPPacket.Info)

type UDPEndpoint interface {
	Close()
	Bind(addr netip.Addr, port int, interfaceId Interface.Id) error
	Listen(funct OnMessageReceivedFunct, errorFunct OnReceiveErrorFunct, appState any) error
	SendTo(addr netip.Addr, port int, handle *system.PacketBufferHandle, interfaceId Interface.Id) error
	SendMsg(pktInfo *IPPacket.Info, msg *system.PacketBufferHandle) error
	LeaveMulticastGroup(interfaceId Interface.Id, addr netip.Addr) error
}
