package responders

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
	"github.com/galenliu/gateway/pkg/dnssd/record"
	"github.com/galenliu/gateway/pkg/matter/inet"
)

type PtrResponder struct {
	*recordResponder
	mTarget *core.FullQName
}

func NewPtrResponder(qName *core.FullQName, target *core.FullQName) *PtrResponder {
	return &PtrResponder{
		recordResponder: &recordResponder{
			responder: &responder{
				mQType: QType.PTR,
				mQName: qName,
			},
			mTtl: kDefaultTtl,
		},
		mTarget: target,
	}
}

func (p *PtrResponder) AddAllResponses(info *inet.IPPacketInfo, delegate ResponderDelegate, configuration *ResponseConfiguration) {
	r := record.NewPtrResourceRecord(p.GetQName(), p.mTarget)
	configuration.Adjust(r)
	delegate.AddResponse(r)
}
