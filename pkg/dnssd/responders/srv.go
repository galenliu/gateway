package responders

import (
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
	"github.com/galenliu/gateway/pkg/dnssd/record"
	"github.com/galenliu/gateway/pkg/matter/inet"
)

type SrvResponder struct {
	*recordResponder
	mRecord *record.SrvResourceRecord
}

func NewSrvResponder(resourceRecord *record.SrvResourceRecord) *SrvResponder {
	return &SrvResponder{
		recordResponder: &recordResponder{
			responder: &responder{
				mQType: QType.SRV,
				mQName: resourceRecord.GetName(),
			},
			mTtl: kDefaultTtl,
		},
		mRecord: resourceRecord,
	}
}

func (r *SrvResponder) AddAllResponses(source *inet.IPPacketInfo, delegate ResponderDelegate, configuration *ResponseConfiguration) {

}
