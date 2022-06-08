package mdns

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"github.com/galenliu/gateway/pkg/matter/inet/udp_endpoint"
)

type ServerDelegate interface {
	OnQuery(data core.BytesRange, info *inet.IPPacketInfo)
	OnResponse(data core.BytesRange, info *inet.IPPacketInfo)
}

type ServerBase interface {
	Shutdown()
	SetDelegate()
	ShutdownEndpoint(info EndpointInfo)
	Listen(manager udp_endpoint.UDPEndpoint, port int)
}
