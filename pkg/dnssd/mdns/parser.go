package mdns

import "github.com/galenliu/gateway/pkg/dnssd/core"

type ParserDelegate interface {
	OnHeader(header *core.ConstHeaderRef)
	OnQuery(queryData *QueryData)
	OnResource(t *core.ResourceType, data *ResourceData)
}

type ResourceData struct {
	mType  core.QType
	mClass core.QClass
	mTtl   uint64
	mData  core.BytesRange
}

type QueryData struct {
	mType              core.QType
	mClass             core.QClass
	mAnswerViaUnicast  bool
	mIsBootAdvertising bool
}

func ParsePacket(packetData *core.BytesRange, delegate ParserDelegate) error {
	return nil
}

func (q *QueryData) Parse(validData *core.BytesRange, start uint8) {
	if validData.Contains(start) {

	}
}
