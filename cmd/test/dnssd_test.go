package test

import (
	"github.com/miekg/dns"
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
