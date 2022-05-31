package dnssd

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/responders"
)

type QueryResponderAllocator struct {
	mAllocatedQNameParts []*responders.RecordResponder
	mQueryResponder      responders.QueryResponder
}

func (q QueryResponderAllocator) Clear() {
	return
}

func (q *QueryResponderAllocator) AllocateQName(serviceType, kCommissionProtocol, kLocalDomain string) core.FullQName {
	return core.FullQName{
		ServerType:         serviceType,
		CommissionProtocol: kCommissionProtocol,
		LocalDomain:        kLocalDomain,
	}
}

func (q *QueryResponderAllocator) AllocateQNameSpace(size uint) {

}

func (q *QueryResponderAllocator) AddAllocatedResponder(res *responders.RecordResponder) responders.QueryResponderSettings {
	q.mAllocatedQNameParts = append(q.mAllocatedQNameParts, res)
	return q.mQueryResponder.AddResponder(res)
}

func (q *QueryResponderAllocator) AddResponder(res *responders.RecordResponder) responders.QueryResponderSettings {
	return q.AddAllocatedResponder(res)
}
