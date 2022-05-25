package endpoint

import (
	"net/netip"
)

type UDPEndPoint struct {
	addr netip.Addr
}

func NewUDPEndPoint() (endpoint *UDPEndPoint) {
	endpoint = &UDPEndPoint{}
	return
}

func (udp UDPEndPoint) Listen() error {
	addr, err := netip.ParseAddr("")
	if err != nil {
		return err
	}
	udp.addr = addr
	return nil
}
