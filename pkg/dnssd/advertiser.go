package dnssd

import (
	"fmt"
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/core/QClass"
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
	"github.com/galenliu/gateway/pkg/dnssd/mdns"
	"github.com/galenliu/gateway/pkg/dnssd/record"
	"github.com/galenliu/gateway/pkg/dnssd/responders"
	"github.com/galenliu/gateway/pkg/errors"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"github.com/galenliu/gateway/pkg/matter/inet/udp_endpoint"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

const kMdnsPort = 5353

type BroadcastAdvertiseType int

const (
	kStarted BroadcastAdvertiseType = iota
	kRemovingAll
)

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

	err := mdns.GlobalServer().StartServer(udpEndPointManager, kMdnsPort)
	if err != nil {
		return err
	}

	err = s.AdvertiseRecords(kStarted)
	if err != nil {
		return err
	}

	s.mIsInitialized = true

	return nil
}

func (s *Advertiser) RemoveServices() error {
	return nil
}

func (s *Advertiser) AdvertiseRecords(t BroadcastAdvertiseType) error {

	responseConfiguration := &responders.ResponseConfiguration{}
	if t == kRemovingAll {
		responseConfiguration.SetTtlSecondsOverride(0)
	}

	//ifs, err := net.Interfaces()
	//if err != nil {
	//	return err
	//}
	//for _, f := range ifs {
	//	if f.Flags < net.FlagMulticast {
	//		continue
	//	}
	//	adders, err := f.Addrs()
	//	if err != nil {
	//		continue
	//	}
	//	for _, a := range adders {
	//		ip, _ := netip.ParseAddr(a.String())
	//		if ip.Is6() {
	//
	//		}
	//	}
	//}
	{
		var packetInfo inet.IPPacketInfo
		packetInfo.SrcPort = kMdnsPort
		packetInfo.DestPort = kMdnsPort
		packetInfo.InterfaceId = inet.InterfaceId{}
		queryData := mdns.NewQueryData(QType.PTR, QClass.IN, false)
		queryData.SetIsInternalBroadcast(true)

		err := s.mResponseSender.Respond(0, queryData, &packetInfo, responseConfiguration)
		if err != nil {
			return err
		}

	}

	//s.mQueryResponderAllocatorCommissionable.GetQueryResponder()
	//s.mQueryResponderAllocatorCommissioner.GetQueryResponder()

	return nil
}

func (s *Advertiser) AdvertiseCommission(params CommissionAdvertisingParameters) error {

	_ = s.AdvertiseRecords(kRemovingAll)

	var allocator *QueryResponderAllocator
	//var serviceType string
	//if params.GetCommissionAdvertiseMode() == KCommissionableNode {
	//	s.mQueryResponderAllocatorCommissionable.Clear()
	//	allocator = s.mQueryResponderAllocatorCommissionable
	//	serviceType = kCommissionableServiceName
	//} else {
	//	s.mQueryResponderAllocatorCommissioner.Clear()
	//	allocator = s.mQueryResponderAllocatorCommissioner
	//	serviceType = kCommissionerServiceName
	//}

	//nameBuffer := s.GetCommissionableInstanceName()
	var serviceType core.ServiceType
	if params.GetCommissionAdvertiseMode() == KCommissionableNode {
		serviceType = kCommissionableServiceName
	} else {
		serviceType = kCommissionableServiceName
	}
	serviceName := &core.FullQName{
		ServerType: serviceType,
		Protocol:   kCommissionProtocol,
		Domain:     kLocalDomain,
	}

	instanceName := &core.FullQName{
		Instance:   s.GetCommissionableInstanceName(),
		ServerType: serviceType,
		Protocol:   kCommissionProtocol,
		Domain:     kLocalDomain,
	}

	hostName := &core.FullQName{
		Domain:   kLocalDomain,
		Instance: s.GetCommissionableInstanceName(),
	}

	if !allocator.AddResponder(responders.NewPtrResponder(serviceName, instanceName)).
		SetReportAdditional(instanceName).
		SetReportInServiceListing(true).
		IsValid() {
		return errors.New("failed to add service PTR record mDNS responder")
	}

	if !allocator.AddResponder(responders.NewSrvResponder(record.NewSrvResourceRecord(instanceName, hostName, params.GetPort()))).
		SetReportAdditional(hostName).
		IsValid() {
		return errors.New("failed to add SRV record mDNS responder")
	}

	if !allocator.AddResponder(responders.NewIPv6Responder(hostName)).
		IsValid() {
		return errors.New("failed to add IPv6 mDNS responder")
	}

	if params.IsIPv4Enabled() {
		if !allocator.AddResponder(responders.NewIPv4Responder(hostName)).
			IsValid() {
			return errors.New("failed to add IPv6 mDNS responder")
		}
	}

	if params.GetVendorId() != nil {
		name := fmt.Sprintf("_V%d", *params.GetVendorId())
		vendorServiceName := &core.FullQName{
			Instance:   name,
			ServerType: kSubtypeServiceNamePart + serviceType,
			Protocol:   kCommissionProtocol,
			Domain:     kLocalDomain,
		}
		if !allocator.AddResponder(responders.NewPtrResponder(vendorServiceName, instanceName)).
			SetReportAdditional(instanceName).
			SetReportInServiceListing(true).
			IsValid() {
			return errors.New("failed to add vendor PTR record mDNS responder")
		}
	}

	if params.GetDeviceType() != nil {
		name := fmt.Sprintf("_T%d", *params.GetDeviceType())
		typeServiceName := &core.FullQName{
			ServerType: kSubtypeServiceNamePart + serviceType,
			Protocol:   kCommissionProtocol,
			Domain:     kLocalDomain,
			Instance:   name,
		}
		if !allocator.AddResponder(responders.NewPtrResponder(typeServiceName, instanceName)).
			SetReportAdditional(instanceName).
			SetReportInServiceListing(true).
			IsValid() {
			return errors.New("failed to add vendor PTR record mDNS responder")
		}
	}

	if params.GetCommissionAdvertiseMode() == KCommissionableNode {
		// TODO
	}

	if !allocator.AddResponder(responders.NewTxtResponder(record.NewTxtResourceRecord(instanceName, s.GetCommissioningTxtEntries(params)))).
		SetReportAdditional(hostName).
		IsValid() {
		return errors.New("failed to add TXT record mDNS responder")
	}

	err := s.AdvertiseRecords(kStarted)
	if err != nil {
		return err
	}
	log.Infof("mDNS service published: %s.%s", instanceName.ServerType, instanceName.Protocol)
	return nil
}

func (s *Advertiser) OnMdnsPacketData(data *core.BytesRange, info *inet.IPPacketInfo) {
	s.mCurrentSource = info
	errors.LogError(mdns.ParsePacket(data, s), "Discovery", "Failed to parse mDNS query")
}

func (s *Advertiser) UpdateCommissionableInstanceName() {
	s.mCommissionableInstanceName = strconv.FormatUint(rand.Uint64(), 16)
	s.mCommissionableInstanceName = strings.ToUpper(s.mCommissionableInstanceName)
}

func (s *Advertiser) GetCommissionableInstanceName() string {
	if s.mCommissionableInstanceName == "" {
		s.mCommissionableInstanceName = strings.Replace(mac48Address(randHex()), ":", "", -1)
	}
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
	var defaultResponseConfiguration = &responders.ResponseConfiguration{}
	errors.LogError(s.mResponseSender.Respond(s.mMessageId, queryData, s.mCurrentSource, defaultResponseConfiguration),
		"Discovery",
		"Failed to reply to query")
}

func (s *Advertiser) OnResource(t *QType.ResourceType, data *mdns.ResourceData) {
	//TODO implement me
	panic("implement me")
}

func (s *Advertiser) FindOperationalAllocator(qname core.QNamePart) *QueryResponderAllocator {
	return nil
}

func (s *Advertiser) GetCommissioningTxtEntries(params CommissionAdvertisingParameters) *core.FullQName {

	var txt = &core.FullQName{
		Txt: make(map[string]string),
	}

	if params.GetProductId() != nil && params.GetVendorId() != nil {
		txt.Txt["VP"] = fmt.Sprintf("%d+%d", *params.GetVendorId(), *params.GetProductId())
	} else if params.GetVendorId() != nil {
		txt.Txt["VP"] = fmt.Sprintf("%d", params.GetVendorId())
	}

	if params.GetDeviceType() != nil {
		txt.Txt["DT"] = fmt.Sprintf("%d", *params.GetDeviceType())
	}

	if params.GetDeviceName() != "" {
		txt.Txt["DN"] = fmt.Sprintf("%s", params.GetDeviceName())
	}
	return txt
}
