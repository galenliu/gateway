package responders

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/core/QClass"
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
	"github.com/galenliu/gateway/pkg/dnssd/record"
	"github.com/galenliu/gateway/pkg/inet/IPPacket"
	"time"
)

type Responder interface {
	ResetAdditionals()
	AddAllResponses(*IPPacket.Info, ResponderDelegate, *ResponseConfiguration)
	GetQName() *core.FullQName
	GetQType() QType.T
	GetQClass() QClass.T
}

type QueryResponderRecord struct {
	Responder
	reportService     bool
	LastMulticastTime time.Time
}

type QueryResponderInfo struct {
	*QueryResponderRecord
	reportNowAsAdditional     bool
	alsoReportAdditionalQName bool
	additionalQName           *core.FullQName
}

type QueryResponderSettings struct {
	mInfo *QueryResponderInfo
}

type QueryResponderRecordFilter struct {
	mIncludeAdditionalRepliesOnly bool
	mReplyFilter                  ReplyFilter
	mIncludeOnlyMulticastBefore   time.Time
}

func (s *QueryResponderSettings) SetReportAdditional(qName *core.FullQName) *QueryResponderSettings {
	if s.IsValid() {
		s.mInfo.alsoReportAdditionalQName = true
		s.mInfo.additionalQName = qName
	}
	return s
}

func (s *QueryResponderSettings) IsValid() bool {
	return s.mInfo != nil
}

func (s *QueryResponderSettings) SetReportInServiceListing(reportService bool) *QueryResponderSettings {
	if s.IsValid() {
		s.mInfo.reportService = reportService
	}
	return s
}

func (f *QueryResponderRecordFilter) SetReplyFilter(filter ReplyFilter) *QueryResponderRecordFilter {
	f.mReplyFilter = filter
	return f
}

func (f *QueryResponderRecordFilter) SetIncludeOnlyMulticastBeforeMS(t time.Time) {
	f.mIncludeOnlyMulticastBefore = t
}

func (f *QueryResponderRecordFilter) Accept(record *QueryResponderInfo) bool {
	if record.Responder == nil {
		return false
	}
	if f.mIncludeAdditionalRepliesOnly && !record.reportNowAsAdditional {
		return false
	}

	if f.mIncludeOnlyMulticastBefore.Before(time.Now()) && record.LastMulticastTime.Before(f.mIncludeOnlyMulticastBefore) {
		return false
	}

	if f.mReplyFilter != nil && !f.mReplyFilter.Accept(record.Responder.GetQType(), record.Responder.GetQClass(), record.Responder.GetQName()) {
		return false
	}

	return true
}

func (f *QueryResponderRecordFilter) SetIncludeAdditionalRepliesOnly(b bool) *QueryResponderRecordFilter {
	f.mIncludeAdditionalRepliesOnly = b
	return f
}

type QueryResponderBase struct {
	*responder
	ResponderInfos []*QueryResponderInfo
}

func (r *QueryResponderBase) ResetAdditionals() {
	for _, r := range r.ResponderInfos {
		r.reportNowAsAdditional = false
	}
}

func (r *QueryResponderBase) AddAllResponses(source *IPPacket.Info, delegate ResponderDelegate, configuration *ResponseConfiguration) {
	for _, m := range r.ResponderInfos {
		if !m.reportService {
			continue
		}
		if m.Responder == nil {
			continue
		}
		r := record.NewPtrResourceRecord(r.GetQName(), m.GetQName())
		configuration.Adjust(r)
		delegate.AddResponse(r)
	}
}

func (r *QueryResponderBase) MarkAdditionalRepliesFor(info *QueryResponderInfo) {
	if !info.alsoReportAdditionalQName {
		return
	}
	if r.markAdditional(info.additionalQName) == 0 {
		return
	}
	var keepAdding = true
	for keepAdding {
		keepAdding = false
		var filter = QueryResponderRecordFilter{}
		filter.SetIncludeAdditionalRepliesOnly(true)
		for _, i := range r.ResponderInfos {
			if i.alsoReportAdditionalQName {
				keepAdding = keepAdding || r.markAdditional(i.additionalQName) != 0
			}
		}

	}

}

func (r *QueryResponderBase) markAdditional(name *core.FullQName) int {
	var count = 0
	for _, info := range r.ResponderInfos {
		if info.reportNowAsAdditional {
			continue
		}
		if info.Responder == nil {
			continue
		}
		if info.GetQName() == r.mQName {
			info.reportNowAsAdditional = true
			count++
		}
	}
	return count
}
