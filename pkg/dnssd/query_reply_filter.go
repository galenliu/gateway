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
	mNameByteRange          *core.BytesRange
	responders.ReplyFilter
}

func NewQueryReplyFilter(q *mdns.QueryData) *QueryReplyFilter {
	return &QueryReplyFilter{
		mIgnoreNameMatch:        false,
		mSendingAdditionalItems: false,
		mQueryData:              q,
	}
}

func (f *QueryReplyFilter) Accept(qType QType.T, qClass QClass.T, fName *core.FullQName) bool {
	if !f.acceptableQueryType(qType) {
		return false
	}

	if !f.acceptableQueryClass(qClass) {
		return false
	}
	return f.acceptablePath(fName)
}

func (f *QueryReplyFilter) acceptableQueryType(qType QType.T) bool {
	if f.mSendingAdditionalItems {
		return true
	}
	return (f.mQueryData.GetType() == QType.ANY) || (f.mQueryData.GetType() == qType)
}

func (f *QueryReplyFilter) acceptableQueryClass(qClass QClass.T) bool {
	return (f.mQueryData.GetClass() == QClass.ANY) || (f.mQueryData.GetClass() == qClass)
}

func (f *QueryReplyFilter) acceptablePath(qName *core.FullQName) bool {
	if f.mIgnoreNameMatch || f.mQueryData.IsInternalBroadcast() {
		return true
	}
	return core.NewFullName(f.mNameByteRange).Equal(qName)
}

func (f *QueryReplyFilter) SetIgnoreNameMatch(b bool) *QueryReplyFilter {
	f.mIgnoreNameMatch = b
	return f
}

func (f *QueryReplyFilter) SetSendingAdditionalItems(b bool) *QueryReplyFilter {
	f.mSendingAdditionalItems = b
	return f
}
