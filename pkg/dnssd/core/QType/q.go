package QType

type QType = uint16

const (
	A         QType = 1
	Ns        QType = 2
	CNAME     QType = 5
	SOA       QType = 6
	NULLVALUE QType = 10
	QTypeWKS  QType = 11
	PTR       QType = 12
	HINFO     QType = 13
	MINFO     QType = 14
	MX        QType = 15
	TXT       QType = 16
	ISDN      QType = 20
	AAAA      QType = 28
	SRV       QType = 33
	DNAM      QType = 39
	ANY       QType = 255
)

type ResourceType = uint

const (
	kAnswer     ResourceType = 0
	kAuthority  ResourceType = 1
	kAdditional ResourceType = 2
)
