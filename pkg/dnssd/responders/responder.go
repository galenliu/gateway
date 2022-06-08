package responders

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/record"
)

type ResponderDelegate interface {
	AddResponse(record *record.ResourceRecord)
}

type Responder struct {
	mQType core.QType
	mQName core.FullQName
}
