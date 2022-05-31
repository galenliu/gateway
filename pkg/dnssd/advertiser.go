package dnssd

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/mdns"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"math/rand"
	"strconv"
	"sync"
)

const KMdnsPort = 5353

type MdnsServerBase interface {
	Shutdown()
}

type Advertiser struct {
	mResponseSender                        *ResponseSender
	mCommissionableInstanceName            string
	mIsInitialized                         bool
	mQueryResponderAllocatorCommissionable *QueryResponderAllocator
	mQueryResponderAllocatorCommissioner   *QueryResponderAllocator
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

func (s *Advertiser) Init(udpEndPointManager inet.UDPEndpoint) error {

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
		s.mQueryResponderAllocatorCommissionable.Clear()
	} else {
		s.mQueryResponderAllocatorCommissioner.Clear()
	}
	name := s.GetCommissionableInstanceName()
	var serviceType string
	var allocator *QueryResponderAllocator
	if params.GetCommissionAdvertiseMode() == KCommissionableNode {
		serviceType = kCommissionableServiceName
		allocator = s.mQueryResponderAllocatorCommissionable
	} else {
		serviceType = kCommissionerServiceName
		allocator = s.mQueryResponderAllocatorCommissioner
	}
	serviceName := allocator.AllocateQName(serviceType, kCommissionProtocol, kLocalDomain)

	hostName := core.FullQName{}

	allocator.AddResponder()
	//instanceName := allocator.AllocateQName(name, serviceType, kCommissionProtocol, kLocalDomain)

	return nil
}

func (s *Advertiser) UpdateCommissionableInstanceName() {
	s.mCommissionableInstanceName = strconv.FormatUint(rand.Uint64(), 10)
}

func (s *Advertiser) GetCommissionableInstanceName() string {
	return s.mCommissionableInstanceName
}

func (s *Advertiser) FinalizeServiceUpdate() error {
	return nil
}
