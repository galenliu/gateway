package mdns

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/inet/IP"
	"github.com/galenliu/gateway/pkg/inet/IPPacket"
	"github.com/galenliu/gateway/pkg/inet/Interface"
	"github.com/galenliu/gateway/pkg/inet/udp_endpoint"
	"github.com/galenliu/gateway/pkg/system"
	"net/netip"
)

type ServerDelegate interface {
	OnQuery(data core.BytesRange, info IPPacket.Info)
	OnResponse(data core.BytesRange, info IPPacket.Info)
}

type ServerBase interface {
	Shutdown()
	SetDelegate()
	ShutdownEndpoint(info EndpointInfo)
	Listen(manager udp_endpoint.UDPEndpoint, port int) error
	DirectSend(packet *system.PacketBufferHandle, address IP.Address, port int, id Interface.Id) error
	BroadcastSend(packet *system.PacketBufferHandle, port int, id Interface.Id, addr IP.Address) error
}

type BroadcastSendDelegate interface {
	Accept(info EndpointInfo) *udp_endpoint.UDPEndpoint
}

type ListenSocketPickerDelegate struct {
	BroadcastSendDelegate
}

type PacketDelegate interface {
	OnMdnsPacketData(data *core.BytesRange, info *IPPacket.Info)
}

type InterfaceTypeFilterDelegate struct {
	BroadcastSendDelegate
	mInterface Interface.Id
	mAddress   netip.Addr
	mChild     BroadcastSendDelegate
}

func (d InterfaceTypeFilterDelegate) Accept(info EndpointInfo) *udp_endpoint.UDPEndpoint {
	return d.mChild.Accept(info)
}
