package record

import "github.com/galenliu/gateway/pkg/dnssd/core"

type SrvResourceRecord struct {
	*resourceRecord
	mServerName *core.FullQName
	mPort       uint16
	mPriority   uint16
	mWeight     uint16
}

func NewSrvResourceRecord(qName *core.FullQName, hostName *core.FullQName, port int) *SrvResourceRecord {
	return &SrvResourceRecord{
		resourceRecord: nil,
		mServerName:    nil,
		mPort:          0,
		mPriority:      0,
		mWeight:        0,
	}
}

func (r *SrvResourceRecord) WriteData(writer *core.RecordWriter) bool {
	return true
}
