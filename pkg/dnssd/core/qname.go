package core

type FullQName struct {
	ServerType         string
	CommissionProtocol string
	LocalDomain        string
}

func NewFullName() *FullQName {
	return &FullQName{}
}
