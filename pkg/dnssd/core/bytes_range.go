package core

import (
	"bytes"
	"encoding/binary"
	"github.com/galenliu/gateway/pkg/dnssd/mdns"
)

const kPtrMask = 0xC0
const kMaxValueSize = 63

// BytesRange 原始的数据
type BytesRange struct {
	data             []byte
	mCurrentPosition uint8
	mIsValid         bool
	mValue           []byte
}

func (r *BytesRange) Contains(data []byte) bool {
	return bytes.Contains(r.Bytes(), data)
}

func (r *BytesRange) Size() uint8 {
	return uint8(len(r.data))
}

func (r *BytesRange) Bytes() []byte {
	return r.data
}

func (r *BytesRange) Get8At(at uint8) uint8 {
	if len(r.data) < int(at) {
		return 0
	}
	data := r.data[at]
	return data
}

func (r *BytesRange) Get16At(offset uint8) uint16 {
	return binary.BigEndian.Uint16(r.data[offset : offset+2])
}

func (r *BytesRange) Next() bool {
	if !r.mIsValid {
		return false
	}
	for {
		if r.Size() < r.mCurrentPosition {
			return false
		}
		length := r.Get8At(r.mCurrentPosition)
		if length == 0 {
			return false
		}

		if (length & kPtrMask) == kPtrMask {

		} else {
			if length > kMaxValueSize {
				r.mIsValid = false
				return false
			}
			if r.Size() < r.mCurrentPosition+1+length {
				r.mIsValid = false
				return false
			}
			r.mValue = r.Bytes()[r.mCurrentPosition+1 : length]
			r.mValue[length] = '0'
			r.mCurrentPosition = r.mCurrentPosition + length + 1
			return true
		}
	}
}

func (r *BytesRange) ParseQueryData() []*mdns.QueryData {
	return nil
}

func (r *BytesRange) ParseQueryResourceData() []*mdns.ResourceData {
	return nil
}

func (r *BytesRange) ParseQueryResourceAdditional() []*mdns.ResourceData {
	return nil
}
