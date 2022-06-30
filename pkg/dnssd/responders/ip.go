package responders

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
)

type IPv4Responder struct {
	*recordResponder
	mTarget *core.FullQName
}

type IPv6Responder struct {
	*recordResponder
}

func NewIPv4Responder(qname *core.FullQName) *IPv4Responder {
	return &IPv4Responder{
		//recordResponder: NewRecordResponder(T.A, qname),
		recordResponder: &recordResponder{
			responder: &responder{
				mQType: QType.A,
				mQName: qname,
			},
			mTtl: kDefaultTtl,
		},
	}
}

func NewIPv6Responder(qname *core.FullQName) *IPv6Responder {
	return &IPv6Responder{
		//recordResponder: NewRecordResponder(T.A, qname),
		recordResponder: &recordResponder{
			responder: &responder{
				mQType: QType.AAAA,
				mQName: qname,
			},
			mTtl: kDefaultTtl,
		},
	}
}

func (p *IPv6Responder) AddAllResponses(info *IPPacket.Info, delegate ResponderDelegate, configuration *ResponseConfiguration) {

}

func (p *IPv4Responder) AddAllResponses(info *IPPacket.Info, delegate ResponderDelegate, configuration *ResponseConfiguration) {

}
