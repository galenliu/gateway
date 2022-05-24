package dnssd

import "github.com/galenliu/gateway/pkg/dnssd/mdns"

const (
	CommssionAdvertiseModeCommissionableNode = iota
	CommssionAdvertiseModeCommissioner
)

const (
	CommissioningModeDisabled        = iota // Commissioning Mode is disabled, CM=0 in DNS-SD key/value pairs
	CommissioningModeEnabledBasic           // Basic Commissioning Mode, CM=1 in DNS-SD key/value pairs
	CommissioningModeEnabledEnhanced        // Enhanced Commissioning Mode, CM=2 in DNS-SD key/value pairs
)

type MdnsServerBase interface {
	Shutdown()
}

type ServiceAdvertiser struct {
	mResponseSender ResponseSender
}

func NewServiceAdvertiser() *ServiceAdvertiser {
	return &ServiceAdvertiser{}
}

func (s ServiceAdvertiser) Init() error {

	server := mdns.NewMdnsServer()

	s.mResponseSender.mServer.Shutdown()
	s.mResponseSender.SetServer(server)
	return nil
}
