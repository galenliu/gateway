package dnssd

import (
	"github.com/galenliu/gateway/pkg/dnssd/mdns"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"math/rand"
)

const MdnsPort = 5353

type MdnsServerBase interface {
	Shutdown()
}

type Advertiser struct {
	mResponseSender             *ResponseSender
	mCommissionableInstanceName uint64
	mdnsServer                  *mdns.MdnsServer
	mIsInitialized              bool
}

func NewAdvertiser() *Advertiser {
	return &Advertiser{}
}

func (s *Advertiser) Init(mgr inet.UDPEndpointManager) error {

	s.mdnsServer = mdns.NewMdnsServer()
	s.mdnsServer.Shutdown()

	if s.mIsInitialized {
		s.UpdateCommissionableInstanceName()
	}

	s.mResponseSender = NewResponseSender()
	s.mResponseSender.SetServer(s.mdnsServer)

	err := s.mdnsServer.StartServer(mgr, MdnsPort)
	if err != nil {
		return err
	}

	s.AdvertiseRecords()

	s.mIsInitialized = true

	return nil
}

func (s *Advertiser) RemoveServices() error {
	return nil
}

func (s *Advertiser) AdvertiseRecords() {

}

func (s *Advertiser) Advertise(params CommissionAdvertisingParameters) error {
	if params.GetCommissionAdvertiseMode() == CommissionableNode {

	} else {

	}

	return nil
}

func (s *Advertiser) UpdateCommissionableInstanceName() {
	s.mCommissionableInstanceName = rand.Uint64()
}

func (s *Advertiser) FinalizeServiceUpdate() error {
	return nil
}
