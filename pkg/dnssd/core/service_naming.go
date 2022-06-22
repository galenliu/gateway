package core

type ServiceType string
type Protocol string

const (
	kCommissionableServiceName ServiceType = "_matterc"
	kCommissionerServiceName   ServiceType = "_matterd"
	kOperationalServiceName    ServiceType = "_matter"
	kCommissionProtocol        Protocol    = "_udp"
	kLocalDomain                           = "local"
	kOperationalProtocol       Protocol    = "_tcp"
)
