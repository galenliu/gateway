package responders

import "github.com/galenliu/gateway/pkg/dnssd/core"

type RecordResponder struct {
	*Responder
}

func NewRecordResponder(qType core.QType, qName core.FullQName) *RecordResponder {
	return &RecordResponder{
		&Responder{
			mQName: qName,
			mQType: qType,
		},
	}
}
