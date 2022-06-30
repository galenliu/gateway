package mdns

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/core/QClass"
	"github.com/galenliu/gateway/pkg/dnssd/core/QType"
	"github.com/galenliu/gateway/pkg/dnssd/core/ResourceType"
)

const kSizeBytes uint8 = 12
const kMaxValueSize = 63

type ParserDelegate interface {
	OnHeader(header *core.ConstHeaderRef)
	OnQuery(queryData *QueryData)
	OnResource(t ResourceType.T, data *ResourceData)
}

type ResourceData struct {
	mType  QType.T
	mClass QClass.T
	mTtl   uint64
	mData  core.BytesRange
}

type QueryData struct {
	mType                QType.T
	mClass               QClass.T
	mAnswerViaUnicast    bool
	mIsInternalBroadcast bool
}

func NewQueryData(qType QType.T, kClass QClass.T, unicast bool) *QueryData {
	return &QueryData{
		mType:                qType,
		mClass:               kClass,
		mAnswerViaUnicast:    unicast,
		mIsInternalBroadcast: false,
	}
}

//func (q *QueryData) Parse(validData *core.BytesRange, start, end uint8) bool {
//	// Structure is:
//	//    QNAME
//	//    TYPE
//	//    CLASS (plus a flag for unicast)
//	if validData.Size() < end {
//		return false
//	}
//	data := validData.Bytes()[start:end]
//	q.mType = validData.Get16At(start)
//	return true
//}

func (q *QueryData) SetIsInternalBroadcast(isInternalBroadcast bool) {
	q.mIsInternalBroadcast = isInternalBroadcast
}

func (q *QueryData) GetType() QType.T {
	return q.mType
}

func (q *QueryData) GetClass() QClass.T {
	return q.mClass
}

func (q *QueryData) IsInternalBroadcast() bool {
	return q.mIsInternalBroadcast
}

func (q *QueryData) GetName() core.FullQName {
	return core.FullQName{
		Instance:   "",
		ServerType: "",
		Protocol:   "",
		Domain:     "",
		Txt:        nil,
	}
}

func (q *QueryData) RequestedUnicastAnswer() bool {
	return q.mAnswerViaUnicast
}

func ParsePacket(packetData *core.BytesRange, delegate ParserDelegate) bool {

	if packetData.Size() < core.KSizeBytes {
		return false
	}
	var header = &core.ConstHeaderRef{Data: packetData.Bytes()}

	if !header.GetFlags().IsValidMdns() {
		return false
	}

	// set messageId
	delegate.OnHeader(header)
	{
		queryDataList := packetData.ParseQueryData()
		for _, queryData := range queryDataList {
			delegate.OnQuery(queryData)
		}

		resourceDataList := packetData.ParseQueryResourceData()
		for _, resourceData := range resourceDataList {
			delegate.OnResource(ResourceType.Answer, resourceData)
		}

		resourceAdditionalList := packetData.ParseQueryResourceAdditional()
		for _, resourceData := range resourceAdditionalList {
			delegate.OnResource(ResourceType.Additional, resourceData)
		}

	}

	return true
}
