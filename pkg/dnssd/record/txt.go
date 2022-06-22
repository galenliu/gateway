package record

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
)

const kTxtDefaultTtl = 4500
const kMaxTxtRecordLength = 63

type TxtResourceRecord struct {
	*resourceRecord
	mEntries string
}

func NewTxtResourceRecord(qName *core.FullQName, entries *core.FullQName) *TxtResourceRecord {
	return &TxtResourceRecord{
		resourceRecord: &resourceRecord{
			mTtl:        kTxtDefaultTtl,
			mQType:      QType.TXT,
			mQname:      qName,
			mCacheFlush: false,
		},
	}
}

func (r *TxtResourceRecord) WriteData(writer *core.RecordWriter) bool {
	return true
}
