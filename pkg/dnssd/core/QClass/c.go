package QClass

type QClass = uint16

const (
	kQClassUnicastAnswerFlag uint16 = 0x8000
	kQClassResponseFlushBit  uint16 = 0x8000
	IN                       QClass = 1
	ANY                      QClass = 255
	InUnicast                QClass = IN | kQClassUnicastAnswerFlag
	InFlush                  QClass = IN | kQClassResponseFlushBit
)
