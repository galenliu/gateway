package record

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
)

const kDefaultTtl = 120

type ResourceRecord interface {
	SetTtl(u uint32)
	WriteData(writer *core.RecordWriter) bool
}

type resourceRecord struct {
	mTtl        uint32
	mQType      QType.QType
	mQname      *core.FullQName
	mCacheFlush bool
}

func NewResourceRecord() *resourceRecord {
	return &resourceRecord{
		mTtl:        kDefaultTtl,
		mQname:      nil,
		mCacheFlush: false,
	}
}

func (r *resourceRecord) SetTtl(u uint32) {
	r.mTtl = u
}

func (r *resourceRecord) setCacheFlush(set bool) {
	r.mCacheFlush = set
}

func (r *resourceRecord) getCacheFlush() bool {
	return r.mCacheFlush
}

func (r *resourceRecord) getTtl() uint32 {
	return r.mTtl
}

func (r *resourceRecord) GetName() *core.FullQName {
	return r.mQname
}

func (r *resourceRecord) append(hdr *core.HeaderRef, asType QType.ResourceType, out *core.RecordWriter) bool {
	return true
}
