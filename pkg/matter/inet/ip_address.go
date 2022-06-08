package inet

type IPAddressType uint8

const (
	Unknown = iota
	IPV4
	IPV6
	Any
)
