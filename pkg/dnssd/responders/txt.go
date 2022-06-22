package responders

import (
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
	"github.com/galenliu/gateway/pkg/dnssd/record"
	"github.com/galenliu/gateway/pkg/matter/inet"
)

type TxtResponder struct {
	*recordResponder
	mRecord *record.TxtResourceRecord
}

func NewTxtResponder(resourceRecord *record.TxtResourceRecord) *TxtResponder {
	return &TxtResponder{
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

func (r *TxtResponder) AddAllResponses(source *inet.IPPacketInfo, delegate ResponderDelegate, configuration *ResponseConfiguration) {

}
