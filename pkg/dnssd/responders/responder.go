package responders

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/core/QClass"
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
	"github.com/galenliu/gateway/pkg/dnssd/record"
)

type ResponderDelegate interface {
	AddResponse(record record.ResourceRecord)
}

type responder struct {
	mQType QType.QType
	mQName *core.FullQName
}

func newResponder(qType QType.QType, mQname *core.FullQName) *responder {
	return &responder{
		mQType: qType,
		mQName: mQname,
	}
}

func (r responder) GetQClass() QClass.QClass {
	return QClass.IN
}

func (r responder) GetQType() QType.QType {
	return r.mQType
}

func (r responder) GetQName() *core.FullQName {
	return r.mQName
}
