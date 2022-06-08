package dnssd

import (
	"github.com/galenliu/gateway/pkg/dnssd/mdns"
	"github.com/galenliu/gateway/pkg/dnssd/responders"
	"github.com/galenliu/gateway/pkg/matter/inet"
)

type ResponseSendingState struct {
	mQuery        *mdns.QueryData
	mSource       *inet.IPPacketInfo
	mMessageId    uint32
	mResourceType *mdns.ResourceData
}

func (state *ResponseSendingState) Reset(messageId uint32, query *mdns.QueryData, packet *inet.IPPacketInfo) {
	state.mMessageId = messageId
	state.mQuery = query
	state.mSource = packet
}

// ResponseSender 实现 ResponderDelegate接口
type ResponseSender struct {
	mServer    mdns.ServerBase
	mSendState *ResponseSendingState
}

func NewResponseSender() *ResponseSender {
	r := &ResponseSender{}
	r.mSendState = &ResponseSendingState{}
	return r
}

func (r ResponseSender) SetServer(ser MdnsServerBase) {
	r.mServer = ser
}

func (r ResponseSender) AddQueryResponder(responder *responders.QueryResponder) {

}

func (r *ResponseSender) Respond(messageId uint32, query *mdns.QueryData, querySource *inet.IPPacketInfo) error {
	r.mSendState.Reset(messageId, query, querySource)

	//TODO implement me
	panic("implement me")
}
