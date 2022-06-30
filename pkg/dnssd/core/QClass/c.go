package QClass

type T = uint16

const (
	kQClassUnicastAnswerFlag uint16 = 0x8000
	kQClassResponseFlushBit  uint16 = 0x8000
	IN                       T      = 1
	ANY                      T      = 255
	InUnicast                T      = IN | kQClassUnicastAnswerFlag
	InFlush                  T      = IN | kQClassResponseFlushBit
)
