package dnssd

import "net/netip"

type EndpointInfo struct {
	network    netip.Addr
	enableIPV4 bool
}

func (info EndpointInfo) IsIpv6() bool {
	return info.network.Is6()
}
