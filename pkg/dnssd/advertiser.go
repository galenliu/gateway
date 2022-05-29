package dnssd

import (
	"github.com/galenliu/gateway/pkg/dnssd/mdns"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"math/rand"
	"sync"
)

const KMdnsPort = 5353

type MdnsServerBase interface {
	Shutdown()
}

type Advertiser struct {
	mResponseSender             *ResponseSender
	mCommissionableInstanceName uint64
	mIsInitialized              bool
}

var insAdvertiser *Advertiser
var advertiserOnce sync.Once

func AdvertiserInstance() *Advertiser {
	advertiserOnce.Do(func() {
		insAdvertiser = newAdvertiser()
	})
	return insAdvertiser
}

func newAdvertiser() *Advertiser {
	return &Advertiser{}
}

func (s *Advertiser) Init(udpEndPointManager inet.UDPEndpointManager) error {

	mdns.GlobalMdnsServer().Shutdown()

	if s.mIsInitialized {
		s.UpdateCommissionableInstanceName()
	}

	s.mResponseSender = NewResponseSender()
	s.mResponseSender.SetServer(mdns.GlobalMdnsServer())

	err := mdns.GlobalMdnsServer().StartServer(udpEndPointManager, KMdnsPort)
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
	if params.GetCommissionAdvertiseMode() == KCommissionableNode {

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
