package record

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
)

type PtrResourceRecord struct {
	*resourceRecord
	mPtrName *core.FullQName
}

func NewPtrResourceRecord(qName, ptrName *core.FullQName) *PtrResourceRecord {
	return &PtrResourceRecord{
		resourceRecord: &resourceRecord{
			mQType:      QType.PTR,
			mQname:      qName,
			mCacheFlush: false,
		},
		mPtrName: ptrName,
	}
}

func (r *PtrResourceRecord) GetPtr() *core.FullQName {
	return r.mPtrName
}

func (r *PtrResourceRecord) WriteData(writer *core.RecordWriter) bool {
	return true
}
