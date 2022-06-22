package core

type QNamePart = []byte

type FullQName struct {
	Instance   string
	ServerType ServiceType
	Protocol   Protocol
	Domain     string
	Txt        map[string]string
}
