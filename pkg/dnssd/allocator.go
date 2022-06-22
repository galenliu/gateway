package dnssd

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/responders"
)

type QueryResponderAllocator struct {
	mAllocatedQNameParts []responders.RecordResponder
	mQueryResponder      *responders.QueryResponderBase
	mFullName            *core.FullQName
}

func (q QueryResponderAllocator) Clear() {
	return
}

func (q *QueryResponderAllocator) AllocateHostQName(mac string, kLocalDomain string) *core.FullQName {
	return &core.FullQName{
		Mac:    mac,
		Domain: kLocalDomain,
	}
}

func (q *QueryResponderAllocator) AllocateQName(serviceType core.ServiceType, kCommissionProtocol core.Protocol, kLocalDomain string, instanceName ...string) *core.FullQName {
	fName := &core.FullQName{
		ServerType: serviceType,
		Protocol:   kCommissionProtocol,
		Domain:     kLocalDomain,
	}
	if instanceName != nil {
		fName.Instance = instanceName[0]
	}
	return fName
}

func (q *QueryResponderAllocator) GetQueryResponder() *responders.QueryResponderBase {
	return q.mQueryResponder
}

func (q *QueryResponderAllocator) AllocateQNameSpace(size uint) {

}

func (q *QueryResponderAllocator) AddAllocatedResponder(res responders.RecordResponder) *responders.QueryResponderSettings {
	q.mAllocatedQNameParts = append(q.mAllocatedQNameParts, res)
	//return q.mQueryResponder.AddResponder(res)
	return nil
}

func (q *QueryResponderAllocator) AddResponder(res responders.RecordResponder) *responders.QueryResponderSettings {
	return q.AddAllocatedResponder(res)
}
