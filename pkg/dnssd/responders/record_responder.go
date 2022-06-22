package responders

import "github.com/galenliu/gateway/pkg/matter/inet"

const kDefaultTtl uint = 4500

type RecordResponder interface {
	AddAllResponses(info *inet.IPPacketInfo, delegate ResponderDelegate, configuration *ResponseConfiguration)
}

type recordResponder struct {
	*responder
	mTtl uint
}
