package core

type QType = uint16

const (
	QTypeA         QType = 1
	QTypeNS        QType = 2
	QTypeCNAME     QType = 5
	QTypeSOA       QType = 6
	QTypeNULLVALUE QType = 10
	QTypeWKS       QType = 11
	QTypePTR       QType = 12
	QTypeHINFO     QType = 13
	QTypeMINFO     QType = 14
	QTypeMX        QType = 15
	QTypeTXT       QType = 16
	QTypeISDN      QType = 20
	QTypeAAAA      QType = 28
	QTypeSRV       QType = 33
	QTypeDNAM      QType = 39
	QTypeANY       QType = 255
)

type QClass = uint16

const (
	kQClassUnicastAnswerFlag uint16 = 0x8000
	kQClassResponseFlushBit  uint16 = 0x8000
	QClassIN                 QClass = 1
	QClassANY                QClass = 255
	QClassIN_UNICAST         QClass = QClassIN | kQClassUnicastAnswerFlag
	QClassIN_FLUSH           QClass = QClassIN | kQClassResponseFlushBit
)

type ResourceType = uint

const (
	kAnswer     ResourceType = 0
	kAuthority  ResourceType = 1
	kAdditional ResourceType = 2
)
