package responders

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/record"
	"github.com/galenliu/gateway/pkg/matter/inet"
	"time"
)

type Responder interface {
	ResetAdditionals()
	AddAllResponses(*inet.IPPacketInfo, ResponderDelegate, *ResponseConfiguration)
	GetQName() *core.FullQName
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

type QueryResponderRecordFilter struct {
	mIncludeAdditionalRepliesOnly bool
	mReplyFilter                  ReplyFilter
	mIncludeOnlyMulticastBefore   time.Time
}

func (f *QueryResponderRecordFilter) SetReplyFilter(filter ReplyFilter) {
	f.mReplyFilter = filter
}

func (f *QueryResponderRecordFilter) SetIncludeOnlyMulticastBeforeMS(t time.Time) {
	f.mIncludeOnlyMulticastBefore = t
}

func (f *QueryResponderRecordFilter) Accept(r *QueryResponderInfo) bool {
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

func (r *QueryResponderBase) AddAllResponses(source *inet.IPPacketInfo, delegate ResponderDelegate, configuration *ResponseConfiguration) {
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
