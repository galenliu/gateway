package mdns

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/core/QClass"
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
)

type ParserDelegate interface {
	OnHeader(header *core.ConstHeaderRef)
	OnQuery(queryData *QueryData)
	OnResource(t *QType.ResourceType, data *ResourceData)
}

type ResourceData struct {
	mType  QType.QType
	mClass QClass.QClass
	mTtl   uint64
	mData  core.BytesRange
}

type QueryData struct {
	mType                QType.QType
	mClass               QClass.QClass
	mAnswerViaUnicast    bool
	mIsInternalBroadcast bool
}

func NewQueryData(qType QType.QType, kClass QClass.QClass, unicast bool) *QueryData {
	return &QueryData{
		mType:                qType,
		mClass:               kClass,
		mAnswerViaUnicast:    unicast,
		mIsInternalBroadcast: false,
	}
}

func ParsePacket(packetData *core.BytesRange, delegate ParserDelegate) error {
	return nil
}

func (q *QueryData) Parse(validData *core.BytesRange, start uint8) {
	if validData.Contains(start) {

	}
}

func (q *QueryData) SetIsInternalBroadcast(isInternalBroadcast bool) {
	q.mIsInternalBroadcast = isInternalBroadcast
}

func (q *QueryData) GetType() QType.QType {
	return q.mType
}

func (q *QueryData) GetClass() QClass.QClass {
	return q.mClass
}

func (q *QueryData) IsInternalBroadcast() bool {
	return q.mIsInternalBroadcast
}

func (q *QueryData) GetName() string {
	return ""
}

func (q *QueryData) RequestedUnicastAnswer() bool {
	return q.mAnswerViaUnicast
}
