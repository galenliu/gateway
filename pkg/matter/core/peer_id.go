package core

type CompressedFabricId uint64
type FabricId uint64

const kUndefinedCompressedFabricId CompressedFabricId = 0

const kUndefinedFabricId FabricId = 0

type PeerId struct {
	mNodeId             NodeId
	mCompressedFabricId CompressedFabricId
}

func NewPeerId(compressedFabricId CompressedFabricId, nodeId NodeId) *PeerId {
	return &PeerId{
		mNodeId:             nodeId,
		mCompressedFabricId: compressedFabricId,
	}
}

func (p PeerId) SetNodeId(id uint64) {
	p.mNodeId = NodeId(id)
}

func (p PeerId) GetNodeId() NodeId {
	return p.mNodeId
}

func (p PeerId) GetCompressedFabricId() CompressedFabricId {
	return p.mCompressedFabricId
}

func (p PeerId) SetCompressedFabricId() NodeId {
	return p.mNodeId
}

func IsValidFabricId(aFabricId FabricId) bool {
	return aFabricId != kUndefinedFabricId
}
