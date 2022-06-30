package dnssd

import (
	"github.com/galenliu/gateway/pkg/dnssd/core/ResourceType"
	"github.com/galenliu/gateway/pkg/dnssd/mdns"
	"github.com/galenliu/gateway/pkg/dnssd/record"
	"github.com/galenliu/gateway/pkg/dnssd/responders"
	"github.com/galenliu/gateway/pkg/inet/IP"
	"github.com/galenliu/gateway/pkg/inet/IPPacket"
	"github.com/galenliu/gateway/pkg/inet/Interface"
	log "github.com/sirupsen/logrus"
	"time"
)

const kMdnsStandardPort = 5353

type ResponseSendingState struct {
	mQuery        *mdns.QueryData
	mSource       *IPPacket.Info
	mMessageId    uint16
	mResourceType ResourceType.T
	mSendError    error
}

func (state *ResponseSendingState) Reset(messageId uint16, query *mdns.QueryData, packet *IPPacket.Info) {
	state.mMessageId = messageId
	state.mQuery = query
	state.mSource = packet
}

func (s *ResponseSendingState) SendUnicast() bool {
	return s.mQuery.RequestedUnicastAnswer() || s.mSource.SrcPort != kMdnsStandardPort
}

func (state *ResponseSendingState) GetError() error {
	return state.mSendError
}

func (state *ResponseSendingState) GetSourcePort() int {
	return state.mSource.SrcPort
}

func (state *ResponseSendingState) GetSourceAddress() IP.Address {
	return state.mSource.SrcAddress
}

func (state *ResponseSendingState) GetSourceInterfaceId() Interface.Id {
	return state.mSource.InterfaceId
}

func (state *ResponseSendingState) SetResourceType(additional ResourceType.T) {
	state.mResourceType = additional
}

// ResponseSender 实现 ResponderDelegate接口
type ResponseSender struct {
	mServer          mdns.ServerBase
	mSendState       *ResponseSendingState
	mResponders      []*responders.QueryResponderBase
	mResponseBuilder ResponseBuilder
}

func NewResponseSender(server mdns.ServerBase) *ResponseSender {
	r := &ResponseSender{mServer: server}
	r.mSendState = &ResponseSendingState{}
	return r
}

func (r *ResponseSender) AddResponse(record record.ResourceRecord) {
	//TODO implement me
	panic("implement me")
}

func (r *ResponseSender) SetServer(ser mdns.ServerBase) {
	r.mServer = ser
}

func (r *ResponseSender) AddQueryResponder(responder *responders.QueryResponderBase) {
	r.mResponders = append(r.mResponders, responder)
}

func (r *ResponseSender) Respond(messageId uint16, query *mdns.QueryData, querySource *IPPacket.Info, configuration *responders.ResponseConfiguration) error {

	r.mSendState.Reset(messageId, query, querySource)

	for _, it := range r.mResponders {
		it.ResetAdditionals()
	}

	// send all 'Answer' replies
	{
		queryReplyFilter := NewQueryReplyFilter(query)
		responseFilter := responders.QueryResponderRecordFilter{}

		responseFilter.SetReplyFilter(queryReplyFilter)

		if !r.mSendState.SendUnicast() {
			responseFilter.SetIncludeOnlyMulticastBeforeMS(time.Now())
		}

		for _, responder := range r.mResponders {
			for _, info := range responder.ResponderInfos {
				if !responseFilter.Accept(info) {
					continue
				}
				info.Responder.AddAllResponses(querySource, r, configuration)
				if err := r.mSendState.GetError(); err != nil {
					return err
				}
				responder.MarkAdditionalRepliesFor(info)
				if !r.mSendState.SendUnicast() {
					info.LastMulticastTime = time.Now()
				}
			}

		}
	}

	// send all 'Additional' replies
	{
		r.mSendState.SetResourceType(ResourceType.Additional)
		queryReplyFilter := NewQueryReplyFilter(query)
		queryReplyFilter.SetIgnoreNameMatch(true).
			SetSendingAdditionalItems(true)

		responseFilter := responders.QueryResponderRecordFilter{}
		responseFilter.SetReplyFilter(queryReplyFilter).
			SetIncludeAdditionalRepliesOnly(true)

		for _, responder := range r.mResponders {
			for _, info := range responder.ResponderInfos {
				if !responseFilter.Accept(info) {
					continue
				}
				info.Responder.AddAllResponses(querySource, r, configuration)
				if err := r.mSendState.GetError(); err != nil {
					return err
				}
			}
		}
	}
	return r.FlushReply()
}

func (r *ResponseSender) FlushReply() error {
	if r.mResponseBuilder.HasResponseRecords() {
		if r.mSendState.SendUnicast() {
			log.Info("Discovery: Directly sending mDns reply to peer %s on port %d", r.mSendState.GetSourcePort())
			err := r.mServer.DirectSend(
				r.mResponseBuilder.ReleasePacket(),
				r.mSendState.GetSourceAddress(),
				r.mSendState.GetSourcePort(),
				r.mSendState.GetSourceInterfaceId())
			if err != nil {
				return err
			}
		} else {
			err := r.mServer.BroadcastSend(
				r.mResponseBuilder.ReleasePacket(),
				kMdnsStandardPort,
				r.mSendState.GetSourceInterfaceId(),
				r.mSendState.GetSourceAddress())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
