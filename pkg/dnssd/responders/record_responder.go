package responders

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
	"github.com/galenliu/gateway/pkg/inet/IPPacket"
)

const kDefaultTtl uint = 4500

type RecordResponder interface {
	AddAllResponses(info *IPPacket.Info, delegate ResponderDelegate, configuration *ResponseConfiguration)
}

type recordResponder struct {
	*responder
}

func newRecordResponder(qType QType.T, qName *core.FullQName) *recordResponder {
	return &recordResponder{
		responder: newResponder(qType, qName),
	}
}
