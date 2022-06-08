package core

type FullQName struct {
	ServerType         string
	CommissionProtocol string
	LocalDomain        string
}

type SerializedQNameIterator struct {
	mValidData       *BytesRange
	mCurrentPosition uint8
	mLookBehindMax   uint8
}

func NewFullName() *FullQName {
	return &FullQName{}
}

func NewSerializedQNameIterator(validData *BytesRange, position uint8) *SerializedQNameIterator {
	return &SerializedQNameIterator{
		mValidData:       validData,
		mCurrentPosition: position,
		mLookBehindMax:   position - validData.Start(),
	}
}
