package dnssd

import (
	"github.com/galenliu/gateway/pkg/dnssd/core"
	"github.com/galenliu/gateway/pkg/dnssd/lib"
	"github.com/galenliu/gateway/pkg/system"
)

type ResponseBuilder struct {
	mPacket       *system.PacketBufferHandle
	mEndianOutput *lib.BufferWriter
	mHeader       *core.HeaderRef
	mWriter       core.RecordWriter
	mBuildOk      bool
}

func (b *ResponseBuilder) HasResponseRecords() bool {
	return b.mHeader.GetAnswerCount() != 0 ||
		b.mHeader.GetAuthorityCount() != 0 ||
		b.mHeader.GetAdditionalCount() != 0
}

func (b *ResponseBuilder) ReleasePacket() *system.PacketBufferHandle {
	b.mHeader = &core.HeaderRef{}
	b.mBuildOk = false
	return b.mPacket
}
