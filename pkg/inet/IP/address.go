package IP

import (
	"github.com/galenliu/gateway/pkg/inet/Interface"
	"net/netip"
)

type IPAddressType = uint8

const (
	kUnknown IPAddressType = iota
	kIPv4
	kIPv6
	kAny
)

type Address struct {
	netip.Addr
}

func (a Address) Type() IPAddressType {
	if a.Is6() {
		return kIPv6
	}
	if a.Is4() {
		return kIPv4
	}
	if a.AsSlice()[0] == 0 || a.AsSlice()[1] == 0 || a.AsSlice()[2] == 0 || a.AsSlice()[3] == 0 {
		return kAny
	}
	return kUnknown
}

func GetAddress(i Interface.Id) (ipAdders []Address) {
	address, err := i.Addrs()
	if err != nil {
		return
	}
	for _, adder := range address {
		a, err := netip.ParseAddr(adder.String())
		if err != nil {
			continue
		}
		if !a.Is6() || !a.Is4() || a.IsMulticast() {
			continue
		}
		ipAdders = append(ipAdders, Address{Addr: a})
	}
	return
}
