package core

import "encoding/binary"

const KSizeBytes = 12

const (
	kIsResponseMask uint16 = 0x8000
	kOpcodeMask     uint16 = 0x7000
	kTruncationMask uint16 = 0x0400
	kReturnCodeMask uint16 = 0x000F
)

const (
	kMessageIdOffset       uint16 = 0
	kFlagsOffset           uint16 = 2
	kQueryCountOffset      uint16 = 4
	kAnswerCountOffset     uint16 = 6
	kAuthorityCountOffset  uint16 = 8
	kAdditionalCountOffset uint16 = 10
)

// BitPackedFlags 数据的标标志位 Flag
type BitPackedFlags struct {
	mValue uint16
}

func (b *BitPackedFlags) IsValidMdns() bool {
	return (b.mValue & (kOpcodeMask | kReturnCodeMask)) == 0
}

/**
 * Allows operations on a DNS header. A DNS Header is defined in RFC 1035
 * and looks like this:
 *
 * | 0| 1 2 3 4| 5| 6| 7| 8| 9| 0| 1| 2 3 4 5 |
 * |               Message ID                 |
 * |QR| OPCODE |AA|TC|RD|RA| Z|AD|CD| RCODE   |
 * |       Items in QUESTION Section          |
 * |       Items in ANSWER Section            |
 * |       Items in AUTHORITY Section         |
 * |       Items in ADDITIONAL Section        |
 */

type ConstHeaderRef struct {
	Data []byte
}

func (r ConstHeaderRef) GetMessageId() uint16 {
	return r.Get16At(kMessageIdOffset)
}

func (h *ConstHeaderRef) GetQueryCount() uint16 {
	return h.Get16At(kQueryCountOffset)
}

func (h *ConstHeaderRef) GetAnswerCount() uint16 {
	return h.Get16At(kAnswerCountOffset)
}

func (h *ConstHeaderRef) GetAuthorityCount() uint16 {
	return h.Get16At(kAuthorityCountOffset)
}

func (h *ConstHeaderRef) GetAdditionalCount() uint16 {
	return h.Get16At(kAdditionalCountOffset)
}

func (h *ConstHeaderRef) GetFlags() *BitPackedFlags {
	return &BitPackedFlags{mValue: h.Get16At(kFlagsOffset)}
}

func (h *ConstHeaderRef) Get16At(offset uint16) uint16 {
	return binary.BigEndian.Uint16(h.Data[offset : offset+2])
}

type HeaderRef struct {
	*ConstHeaderRef
}
