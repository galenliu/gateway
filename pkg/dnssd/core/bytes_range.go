package core

type BytesRange struct {
	mStart uint8
	mEnd   uint8
}

func NewByteRange(start, end uint8) *BytesRange {
	return &BytesRange{
		mStart: start,
		mEnd:   end,
	}
}

func (r *BytesRange) Contains(start uint8) bool {
	return r.mStart < start && r.mEnd < start
}

func (r *BytesRange) Start() uint8 {
	return r.mStart
}
