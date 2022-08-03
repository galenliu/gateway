package test

import (
	"fmt"
	"github.com/miekg/dns"
	"math/rand"
	"net/netip"
	"testing"
)

func TestName(t *testing.T) {

	t.Log(ParseFullQName("gg", "cc", "local"))
	m := dns.Msg{
		Question: make([]dns.Question, 1),
	}

	m.Question[0] = dns.Question{
		Name:   "pi.local.",
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}
	t.Log(m.MsgHdr.String())
	t.Log("----------------")
	//t.Log(m.MsgHdr.String())
	t.Log(m.String())
	t.Log("----------------")
	t.Log(m.Question[0].String())
}

func ParseFullQName(args ...string) string {
	var s string
	for _, a := range args {
		s = s + dns.Fqdn(a)
	}
	s = dns.Fqdn(s)
	return s
}

func Multicast(id uint64, gid uint16) netip.Addr {
	var scope uint8 = 0x05

	//lFlagsAndScope = (Ox01 << 4  =  0x10)  | 050  = 0x15
	var lFlagsAndScope uint8 = ((01 & 0xF) << 4) | (scope & 0xF)
	var lReserved uint8 = 0x0
	var prefixLength uint8 = 0x40

	var prefix uint64 = 0xfd00000000000000 | (uint64(id) >> 8 & 0x00ffffffffffffff)

	var groupId uint32 = ((uint32(id) << 24) & 0xff000000) | uint32(gid)

	ipV6 := netip.AddrFrom16([16]byte{
		0xFF, lFlagsAndScope, lReserved, prefixLength,
		byte(prefix >> 56), byte(prefix >> 48), byte(prefix >> 40), byte(prefix >> 32),
		byte(prefix >> 24), byte(prefix >> 16), byte(prefix >> 8), byte(prefix >> 0),
		byte(groupId >> 24), byte(groupId >> 16), byte(groupId >> 8), byte(groupId >> 0),
	})
	return ipV6
}

func TestSliceUint64(t *testing.T) {
	u64 := rand.Uint64()
	u16 := uint16(rand.Uint32())
	fmt.Printf("U64: %016x \t\n", u64)
	fmt.Printf("U16: %04x \t\n", u16)
	ipV6 := Multicast(u64, u16)

	fmt.Printf("IPv6: %s \t\n", ipV6.String())
	fmt.Printf("IPv6 IsLinkLocalUnicast 本地链路地址: %t \t\n", ipV6.IsLinkLocalUnicast())
	fmt.Printf("IPv6 IsGlobalUnicast 全球单播地址: %t \t\n", ipV6.IsGlobalUnicast())
	fmt.Printf("IPv6 IsMulticast 组播地址: %t \t\n", ipV6.IsMulticast())

}
