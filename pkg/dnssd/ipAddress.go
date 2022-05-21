package dnssd

type IPAddressType = int

const (
	IPTUnknown IPAddressType = iota
	IPV4
	IPV6
)

func GetIpv6Into() string {
	return ""
}
