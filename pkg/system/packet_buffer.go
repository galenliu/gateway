package system

const kStructureSize = 4

type buf struct {
}

type PacketBuffer struct {
	next   *buf
	totLen uint16
	len    uint16
	ref    uint16
}
