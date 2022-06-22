package dnssd

import (
	"github.com/galenliu/gateway/pkg/dnssd/mdns"
	"github.com/galenliu/gateway/pkg/dnssd/record"
	"github.com/galenliu/gateway/pkg/dnssd/responders"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"time"
)

const kMdnsStandardPort = 5353

type ResponseSendingState struct {
	mQuery        *mdns.QueryData
	mSource       *inet.IPPacketInfo
	mMessageId    uint32
	mResourceType *mdns.ResourceData
	mSendError    error
}

func (state *ResponseSendingState) Reset(messageId uint32, query *mdns.QueryData, packet *inet.IPPacketInfo) {
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

// ResponseSender 实现 ResponderDelegate接口
type ResponseSender struct {
	mServer     mdns.ServerBase
	mSendState  *ResponseSendingState
	mResponders []responders.QueryResponderBase
}

func (r *ResponseSender) AddResponse(record record.ResourceRecord) {
	//TODO implement me
	panic("implement me")
}

func NewResponseSender() *ResponseSender {
	r := &ResponseSender{}
	r.mSendState = &ResponseSendingState{}
	return r
}

func (r ResponseSender) SetServer(ser MdnsServerBase) {

}

func (r ResponseSender) AddQueryResponder(responder *responders.QueryResponderBase) {

}

func (r *ResponseSender) Respond(messageId uint32, query *mdns.QueryData, querySource *inet.IPPacketInfo, configuration *responders.ResponseConfiguration) error {

	r.mSendState.Reset(messageId, query, querySource)

	for _, it := range r.mResponders {
		it.ResetAdditionals()
	}

	// send all 'Answer' replies
	{
		queryReplyFilter := NewQueryReplyFilter(query)
		responseFilter := responders.QueryResponderRecordFilter{}
		responseFilter.SetReplyFilter(queryReplyFilter)
		if r.mSendState.SendUnicast() {
			responseFilter.SetIncludeOnlyMulticastBeforeMS(time.Now())
		}

		for _, responder := range r.mResponders {
			for _, info := range responder.ResponderInfos {
				if !responseFilter.Accept(info) {
					continue
				}
				info.AddAllResponses(querySource, r, configuration)
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
	return r.FlushReply()
}

func (r *ResponseSender) FlushReply() error {
	return nil
}
