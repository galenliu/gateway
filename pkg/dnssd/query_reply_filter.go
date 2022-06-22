package dnssd

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/core/QClass"
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
	"github.com/galenliu/gateway/pkg/dnssd/mdns"
	"github.com/galenliu/gateway/pkg/dnssd/responders"
)

type QueryReplyFilter struct {
	mIgnoreNameMatch        bool
	mSendingAdditionalItems bool
	mQueryData              *mdns.QueryData
	responders.ReplyFilter
}

func NewQueryReplyFilter(q *mdns.QueryData) *QueryReplyFilter {
	return &QueryReplyFilter{
		mIgnoreNameMatch:        false,
		mSendingAdditionalItems: false,
		mQueryData:              q,
	}
}

func (f *QueryReplyFilter) Accept(qType QType.QType, qClass QClass.QClass, fName *core.FullQName) bool {
	if !f.acceptableQueryType(qType) {
		return false
	}

	if !f.acceptableQueryClass(qClass) {
		return false
	}
	return f.acceptablePath(fName)
}

func (f *QueryReplyFilter) acceptableQueryType(qType QType.QType) bool {
	if f.mSendingAdditionalItems {
		return true
	}
	return (f.mQueryData.GetType() == QType.ANY) || (f.mQueryData.GetType() == qType)
}

func (f *QueryReplyFilter) acceptableQueryClass(qClass QClass.QClass) bool {
	return (f.mQueryData.GetClass() == QClass.ANY) || (f.mQueryData.GetClass() == qClass)
}

func (f *QueryReplyFilter) acceptablePath(qName *core.FullQName) bool {
	if f.mIgnoreNameMatch || f.mQueryData.IsInternalBroadcast() {
		return true
	}

	// ？？？
	return f.mQueryData.GetName() == qName.Instance
}
