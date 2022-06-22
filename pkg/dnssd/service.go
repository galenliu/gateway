package dnssd

import (
	"fmt"
	"net/netip"
)

const (
	kSubtypeServiceNamePart    = "_sub"
	kCommissionableServiceName = "_matterc"
	kCommissionerServiceName   = "_matterd"
	kOperationalServiceName    = "_matter"
	kCommissionProtocol        = "_udp"
	kLocalDomain               = "local"
	kOperationalProtocol       = "_tcp"
)

type Config struct {
	Name   string
	Type   string
	Domain string
	Host   string
	Text   map[string]string
	IPs    []netip.Addr
	Port   int
	Ifaces []string
}

func NewConf(instanceName string) *Config {
	t := fmt.Sprintf("%s.%s", kCommissionableServiceName, kCommissionProtocol)
	return &Config{
		Name:   instanceName,
		Type:   t,
		Domain: kLocalDomain,
		Host:   "",
		Text:   nil,
		IPs:    nil,
		Port:   0,
		Ifaces: nil,
	}
}
