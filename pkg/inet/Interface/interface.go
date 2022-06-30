package Interface

import (
	"net"
)

type Type = uint

const (
	IUnknown Type = iota
	WiFi
	Ethernet
	Cellular
	Thread
)

type Id struct {
	net.Interface
}

func (i Id) IsPresent() bool {
	return false
}

func GetInterfaceIds() (list []Id) {
	l, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, i := range l {
		list = append(list, Id{i})
	}
	return list
}
