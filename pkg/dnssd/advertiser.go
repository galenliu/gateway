package dnssd

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/mdns"
	"github.com/galenliu/gateway/pkg/errors"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"github.com/galenliu/gateway/pkg/matter/inet/udp_endpoint"
	"math/rand"
	"strconv"
	"sync"
)

const KMdnsPort = 5353

type MdnsServerBase interface {
	Shutdown()
}

// Advertiser 实现 PacketDelegate 和 ParserDelegate
type Advertiser struct {
	mResponseSender                        *ResponseSender
	mCommissionableInstanceName            string
	mIsInitialized                         bool
	mQueryResponderAllocatorCommissionable *QueryResponderAllocator
	mQueryResponderAllocatorCommissioner   *QueryResponderAllocator
	mCurrentSource                         *inet.IPPacketInfo
	mMessageId                             uint32
}

var insAdvertiser *Advertiser
var advertiserOnce sync.Once

func AdvertiserInstance() *Advertiser {
	advertiserOnce.Do(func() {
		insAdvertiser = newAdvertiser()
		mdns.GlobalServer().SetQueryDelegate(insAdvertiser)
		insAdvertiser.mResponseSender.AddQueryResponder(insAdvertiser.mQueryResponderAllocatorCommissionable.GetQueryResponder())
		insAdvertiser.mResponseSender.AddQueryResponder(insAdvertiser.mQueryResponderAllocatorCommissioner.GetQueryResponder())
	})
	return insAdvertiser
}

func newAdvertiser() *Advertiser {
	return &Advertiser{}
}

func (s *Advertiser) Init(udpEndPointManager udp_endpoint.UDPEndpoint) error {

	mdns.GlobalServer().Shutdown()

	if s.mIsInitialized {
		s.UpdateCommissionableInstanceName()
	}

	s.mResponseSender = NewResponseSender()
	s.mResponseSender.SetServer(mdns.GlobalServer())

	err := mdns.GlobalServer().StartServer(udpEndPointManager, KMdnsPort)
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
	//name := s.GetCommissionableInstanceName()
	//var serviceType string
	//var allocator *QueryResponderAllocator
	//if params.GetCommissionAdvertiseMode() == KCommissionableNode {
	//	serviceType = kCommissionableServiceName
	//	allocator = s.mQueryResponderAllocatorCommissionable
	//} else {
	//	serviceType = kCommissionerServiceName
	//	allocator = s.mQueryResponderAllocatorCommissioner
	//}
	//serviceName := allocator.AllocateQName(serviceType, kCommissionProtocol, kLocalDomain)
	//
	//hostName := core.FullQName{}

	//allocator.AddResponder()
	//instanceName := allocator.AllocateQName(name, serviceType, kCommissionProtocol, kLocalDomain)

	return nil
}

func (s *Advertiser) OnMdnsPacketData(data *core.BytesRange, info *inet.IPPacketInfo) {
	s.mCurrentSource = info
	errors.LogError(mdns.ParsePacket(data, s), "Discovery", "Failed to parse mDNS query")
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

func (s *Advertiser) OnHeader(header *core.ConstHeaderRef) {
	//TODO implement me
	panic("implement me")
}

func (s *Advertiser) OnQuery(queryData *mdns.QueryData) {
	if s.mCurrentSource == nil {
		return
	}
	errors.LogError(s.mResponseSender.Respond(s.mMessageId, queryData, s.mCurrentSource),
		"Discovery",
		"Failed to reply to query")
}

func (s *Advertiser) OnResource(t *core.ResourceType, data *mdns.ResourceData) {
	//TODO implement me
	panic("implement me")
}
