package responders

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"time"
)

type QueryResponderRecord struct {
	responder         *Responder
	reportService     bool
	lastMulticastTime time.Time
}

type QueryResponderInfo struct {
	*QueryResponderRecord
	reportNowAsAdditional      bool
	alsoReportAdditionalQNamef bool
	additionalQName            core.FullQName
}

type QueryResponderSettings struct {
	mInfo *QueryResponderInfo
}

type QueryResponderRecordFilter struct {
}

type QueryResponderIterator struct {
}

type QueryResponderBase struct {
	*Responder
	mResponderInfos []*QueryResponderInfo
}

type QueryResponder struct {
	*QueryResponderBase
}

func (r *QueryResponderBase) AddResponder(res *RecordResponder) QueryResponderSettings {
	info := &QueryResponderInfo{}
	info.responder = res.Responder
	r.mResponderInfos = append(r.mResponderInfos, info)

	return QueryResponderSettings{mInfo: info}
}
